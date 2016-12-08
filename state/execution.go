package state

import (
	"bytes"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
	tmsp "github.com/tendermint/tmsp/types"
	"github.com/zballs/comit/forms"
	"github.com/zballs/comit/types"
	. "github.com/zballs/comit/util"
)

// If the tx is invalid, a TMSP error will be returned.
func ExecTx(state *State, tx types.Tx, isCheckTx bool) (res tmsp.Result) {

	chainID := state.GetChainID()

	// Validate Input Basic
	res = tx.Input.ValidateBasic()
	if res.IsErr() {
		return res
	}

	var inAcc *types.Account

	if tx.Type == types.CreateAccountTx {
		// Create new account
		// Must have txIn pubKey
		inAcc = types.NewAccount(tx.Input.PubKey, 0)
	} else {
		// Get input account
		inAcc = state.GetAccount(tx.Input.Address)
		if inAcc == nil {
			return tmsp.ErrBaseUnknownAddress
		}
		if tx.Input.PubKey != nil {
			inAcc.PubKey = tx.Input.PubKey
		}
	}

	// Validate input, advanced
	signBytes := tx.SignBytes(chainID)
	res = validateInputAdvanced(inAcc, signBytes, tx.Input)
	if res.IsErr() {
		log.Info(Fmt("validateInputAdvanced failed on %X: %v", tx.Input.Address, res))
		return res.PrependLog("in validateInputAdvanced()")
	}

	inAcc.Sequence += 1

	// If CheckTx, we are done.
	if isCheckTx {
		state.SetAccount(tx.Input.Address, inAcc)
		return tmsp.OK
	}

	// Create inAcc checkpoint
	inAccCopy := inAcc.Copy()

	// Run the tx.
	cacheState := state.CacheWrap()
	cacheState.SetAccount(tx.Input.Address, inAcc)
	switch tx.Type {
	case types.CreateAccountTx:
		res = RunCreateAccountTx()
	case types.RemoveAccountTx:
		if inAcc.IsAdmin() {
			res = tmsp.ErrUnauthorized
		} else {
			res = RunRemoveAccountTx(cacheState, tx.Input.Address)
		}
	case types.CreateAdminTx:
		if !inAcc.PermissionToCreateAdmin() {
			res = tmsp.ErrUnauthorized
		} else {
			res = RunCreateAdminTx(cacheState, tx.Data)
		}
	case types.RemoveAdminTx:
		if !inAcc.IsAdmin() {
			res = tmsp.ErrUnauthorized
		} else {
			res = RunRemoveAdminTx(cacheState, tx.Input.Address)
		}
	case types.SubmitTx:
		res = RunSubmitTx(cacheState, tx.Data)
	case types.ResolveTx:
		if !inAcc.PermissionToResolve() {
			res = tmsp.ErrUnauthorized
		} else {
			pubKey := inAcc.PubKey.(crypto.PubKeyEd25519)
			res = RunResolveTx(cacheState, pubKey, tx.Data)
		}
	default:
		res = tmsp.ErrUnknownRequest.SetLog(
			Fmt("Error unrecognized tx type: %v", tx.Type))
	}
	if res.IsOK() {
		cacheState.CacheSync()
		log.Info("Successful execution")
	} else {
		log.Info("AppTx failed", "error", res)
		cacheState.SetAccount(tx.Input.Address, inAccCopy)
	}
	return res
}

//=====================================================================//

func RunCreateAccountTx() tmsp.Result {
	// Just return OK
	return tmsp.OK
}

func RunRemoveAccountTx(state *State, address []byte) tmsp.Result {
	// Return key so we can remove in AppendTx
	key := AccountKey(address)
	return tmsp.NewResultOK(key, "")
}

func RunCreateAdminTx(state *State, data []byte) tmsp.Result {

	// Get secret
	secret, _, err := wire.GetByteSlice(data)
	if err != nil {
		return tmsp.ErrEncodingError.SetLog(
			Fmt("Error: could not get secret: %v", data))
	}

	// Create keys
	pubKey, privKey := CreateKeys(secret)

	// Create new admin
	newAcc := types.NewAdmin(pubKey)
	state.SetAccount(pubKey.Address(), newAcc)

	// Return PubKeyBytes
	buf, n, err := new(bytes.Buffer), int(0), error(nil)
	wire.WriteByteSlice(pubKey[:], buf, &n, &err)
	wire.WriteByteSlice(privKey[:], buf, &n, &err)
	return tmsp.NewResultOK(buf.Bytes(), "account")
}

func RunRemoveAdminTx(state *State, address []byte) tmsp.Result {
	// Return key so we can remove in AppendTx
	key := AccountKey(address)
	return tmsp.NewResultOK(key, "account")
}

func RunSubmitTx(state *State, data []byte) (res tmsp.Result) {
	var form forms.Form
	err := wire.ReadBinaryBytes(data, &form)
	if err != nil {
		return tmsp.ErrEncodingError.SetLog(
			Fmt("Error: could not decode form data: %v", data))
	}
	issue := form.Issue
	formID := (&form).ID()
	buf, n, err := new(bytes.Buffer), int(0), error(nil)
	wire.WriteByteSlice(formID, buf, &n, &err)
	state.Set(buf.Bytes(), data)

	err = state.AddToFilter(buf.Bytes(), issue)
	if err != nil {
		// False positive
	}

	formBytes := wire.BinaryBytes(form)
	data = make([]byte, wire.ByteSliceSize(formBytes)+1)
	bz := data
	bz[0] = types.SubmitTx
	bz = bz[1:]
	wire.PutByteSlice(bz, formBytes)

	return tmsp.NewResultOK(data, "")
}

func RunResolveTx(state *State, pubKey crypto.PubKeyEd25519, data []byte) (res tmsp.Result) {
	formID, _, err := wire.GetByteSlice(data)
	if err != nil {
		return tmsp.NewResult(
			forms.ErrDecodingFormID, nil, "")
	}
	value := state.Get(data)
	if len(value) == 0 {
		return tmsp.NewResult(
			forms.ErrFindForm, nil, Fmt("Error cannot find form with ID: %X", formID))
	}
	var form forms.Form
	err = wire.ReadBinaryBytes(value, &form)
	if err != nil {
		return tmsp.ErrEncodingError.SetLog(
			Fmt("Error parsing form bytes: %v", err.Error()))
	}
	minuteString := ToTheMinute(TimeString())
	pubKeyString := BytesToHexString(pubKey[:])
	err = (&form).Resolve(minuteString, pubKeyString)
	if err != nil {
		return tmsp.NewResult(
			forms.ErrFormAlreadyResolved, nil, Fmt("Error already resolved form with ID: %v", formID))
	}
	buf, n, err := new(bytes.Buffer), int(0), error(nil)
	wire.WriteBinary(form, buf, &n, &err)
	if err != nil {
		return tmsp.ErrEncodingError.SetLog(
			Fmt("Error encoding form with ID: %v", formID))
	}
	state.Set(data, buf.Bytes())

	/*
		err = state.AddToFilter(data, "resolved")
		if err != nil {
			// False positive
		}
	*/

	data = make([]byte, wire.ByteSliceSize(buf.Bytes())+1)
	bz := data
	bz[0] = types.ResolveTx
	bz = bz[1:]
	wire.PutByteSlice(bz, buf.Bytes())

	return tmsp.NewResultOK(data, "")
}

//===============================================================================================//

func validateInputAdvanced(acc *types.Account, signBytes []byte, in types.TxInput) (res tmsp.Result) {
	// Check sequence
	seq := acc.Sequence
	if seq+1 != in.Sequence {
		return tmsp.ErrBaseInvalidSequence.AppendLog(
			Fmt("Got %v, expected %v. (acc.seq=%v)", in.Sequence, seq+1, acc.Sequence))
	}
	// Check signatures
	if !acc.PubKey.VerifyBytes(signBytes, in.Signature) {
		return tmsp.ErrBaseInvalidSignature.AppendLog(
			Fmt("SignBytes: %X", signBytes))
	}
	return tmsp.OK
}
