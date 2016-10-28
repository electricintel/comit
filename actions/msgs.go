package actions

const (
	invalid_hex = "<span style='color:red;'>Invalid hex string: <small>%X</small></span>"

	create_account_success = "<strong style='opacity:0.8;'>public key</strong> <really-small>%X</really-small><br><strong style='opacity:0.8;'>private key</strong> <really-small>%X</really-small>"
	create_account_failure = "<span style='color:red;'>Could not create account</span>"

	remove_account_failure = "<span style='color:red;'>Could not remove account<br><strong>public key</strong> <small>%v</small>"
	remove_account_success = "Removed account<br><strong>public key</strong> <small>%X</small>"

	create_admin_failure = "<span style='color:red;'>Failed to create admin</span>"
	create_admin_success = ""

	remove_admin_failure = "<span style='color:red;'>Could not remove admin<br><strong>Public Key</strong> <small>%v</small><br><strong>Passphrase</strong> <small>%v</small></span>"
	remove_admin_success = "Removed admin<br><strong>Public Key</strong> <small>%v</small><br><strong>Passphrase</strong> <small>%v</small>"

	submit_form_failure = "<span style='color:red;'>Failed to submit form</span>"
	submit_form_success = "<strong>Submitted form with ID</strong> <small>%X</small>"

	form_already_exists = "<span style='color:red;'>Form already exists</span>"

	find_form_failure   = "<span style='color:red;'>Failed to find<br><strong>form ID</strong> <small>%X</small></span>"
	decode_form_failure = "<span style='color:red;'>Failed to decode<br><strong>form ID</strong> <small>%v</small></span>"

	resolve_form_failure = "<span style='color:red;'>Failed to resolve<br><strong>form ID</strong> <small>%v</small></span>"
	resolve_form_success = "Resolved<br><strong>Form ID</strong> <small>%v</small>"

	search_forms_failure = "<span style='color:red;'>Failed to find forms</span>"

	already_connected = "<span style='color:red;'>You are already connected to the network</span>"
	connect_failure   = "<span style='color:red;'>Failed to connect to network</span>"

	find_admin_failure   = "<span style='color:red;'>Could not find admin with public key <small>%v</small></span>"
	send_message_failure = "<span style='color:red;'>Failed to send message to <small>%v</small></span>"
	send_message_success = "<span style='color:red;'>Sent message to <small>%v</small></span>"

	select_option = "<option value='%v'>%v</option>"

/*
	unauthorized         = "<span style='color:red;'>Unauthorized</span>"
	find_admin_failure = "<span style='color:red;'>Failed to find admin</span>"
	admin_update = "admin-%v-update"
	calc_metric_success = "<strong style='opacity: 0.9;'>%v</strong><br> <small>%v</small>"
	calc_metric_failure = "<span style='color:red;'>Could not calculate metric -- zero forms found</span>"
*/
)
