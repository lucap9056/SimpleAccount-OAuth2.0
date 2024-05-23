import { writable, type Writable } from "svelte/store";

type Translation = {
    [key: string]: string
}

const _default = {
    "myapps": "My Apps",
    "myapps_new": "New App",
    "app_info": "App Info",
    "confirm": "Confirm",
    "cancel": "Cancel",
    "create": "Create",
    "back": "Back",
    "edit": "Edit",
    "delete": "Delete",

    "edit_app": "Edit App",
    "create_app": "Create App",
    "name": "Name",
    "callback": "Callback",

    "edit_app_name_limit": "Available characters: a-z A-Z 0-9 _-",
    "edit_app_name_too_short": "Your app name must be at least 4 characters long.",
    "edit_app_name_too_long": "App names cannot be more than 32 characters.",
    "edit_invalid_app_name": "App name contains invalid characters.",

    "edit_callback_placeholder": "https://www.example.com/login?code={code}",
    "edit_callback_hint": "Verification code automatically replaced by {code}.",
    "edit_invalid_callback_url": "Invalid URL format.",
    "edit_callback_url_missing_placeholder": "The callback URL is missing the placeholder: {code}.",

    "edit_label_secret": "Save the verification information carefully.",
    "edit_label_token": "Token",
    "edit_label_private_key": "Private Key",
    "edit_value_copy": "Copied",
    "confirm_delete_title": "re you sure you want to delete?",
    "confirm_delete_message": "his action cannot be undone. Do you really want to proceed?",

    "permissions_requested": "The permissions requested by this app:",
    "permission_read_id": "Read ID",
    "permission_read_username": "Read Username",
    "permission_read_email": "Read Email",
    "permission_read_last_edit": "Read Last Edit Time",
    "permission_read_register_time": "Read Register Time",

    "error_server_side": "An error occurred on the server side.",
    "error_not_allowed_to_create_apps": "You are not allowed to create apps.",
    "error_not_logged_in": "You are not logged in.",
    "error_client_invalid_request": "Invalid client request.",
    "error_app_name_is_empty": "The app name is empty.",
    "error_app_callback_is_empty": "The app callback is empty.",
    "error_user_not_exist": "User does not exist.",
    "error_app_not_exist": "The app does not exist.",
    "error_verification_code_invalid": "Invalid verification code.",
    "error_app_name_already_exists": "The app name already exists.",
    "error_owned_apps_reached_limit": "Number of owned apps has reached the limit.",
};

let translation: Translation = _default;
const Translations: Writable<Translation> = writable(translation);
Translations.subscribe((l) => {
    translation = l;
});

fetch(`./translations/${navigator.language}.json`).then(async (res) => {
    if (!res.ok) return;

    let body = "";
    const reader = res.body.getReader();
    const decoder = new TextDecoder();
    while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        body += decoder.decode(value, { stream: true });
    }

    const Translation: Translation = JSON.parse(body);
    Translations.update(() => Translation);
});

export default Translations;

export function TranslationsGet(id: string) {
    return translation[id] || id;
}