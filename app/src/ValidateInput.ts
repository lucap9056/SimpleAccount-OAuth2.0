import Translations, { TranslationsGet } from "./Translations";

function AppName(value: string): string {
    if (!/^[a-zA-Z0-9_-]+$/.test(value)) {
        return TranslationsGet("edit_invalid_app_name");
    }

    if (value.length < 4) {
        return TranslationsGet("edit_app_name_too_short");
    }

    if (value.length > 32) {
        return TranslationsGet("edit_app_name_too_length");
    }

    return "";
}

function AppCallback(value: string): string {
    if (!/^(?:\w+:)?\/\/([^\s.]+\.\S{2}|localhost[\:?\d]*)\S*$/.test(value)) {
        return TranslationsGet("edit_invalid_callback_url");
    }

    if (!/{code}/.test(value)) {
        return TranslationsGet("edit_callback_url_missing_placeholder");
    }
    
    return "";
}

export default {
    AppName,
    AppCallback,
}