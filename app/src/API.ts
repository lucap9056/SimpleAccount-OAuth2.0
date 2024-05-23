import Authorization from "./Authorization";

const ErrorCode: { [key: number]: string } = {
    0: "error_null",
    1: "error_server_side",
    11: "error_not_allowed_to_create_apps",
    12: "error_not_logged_in",
    13: "error_client_invalid_request",
    14: "error_app_name_is_empty",
    15: "error_app_callback_is_empty",
    16: "error_user_not_exist",
    17: "error_app_not_exist",
    18: "error_verification_code_invalid",
    19: "error_app_name_already_exists",
    20: "error_owned_apps_reached_limit",
}

const MainHost = "/api";
const Host = "/oauth_api";

export const Path: { [key: string]: string } = {
    USER: MainHost + "/user/",
    //GET: request user data
    INFO: Host + "/info/",
    //GET: request app data
    //POST: edit app data
    APP: Host + "/app/",
    //GET: app
}

type RawResponseData = {
    success: boolean
    result: string
    error: number
}

export function Post(path: string, postBody?: { [key: string]: string | number | boolean }): Promise<ResponseData> {
    return Request("POST", path, postBody);
}

export function Put(path: string, putBody?: { [key: string]: string | number | boolean }): Promise<ResponseData> {
    return Request("PUT", path, putBody);
}

export function Get(path: string): Promise<ResponseData> {
    return Request("GET", path);
}

function Request(method: string, path: string, body?: { [key: string]: string | number | boolean }): Promise<ResponseData> {
    return new Promise(async (reslove, reject) => {

        const data = {
            method,
            headers: {
                "Content-Type": "application/json",
                "Authorization": Authorization.GetToken()
            },
        }

        if (body !== null) {
            data["body"] = JSON.stringify(body);
        }
        const response = await fetch(path, data);

        let responseBody = "";
        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            responseBody += decoder.decode(value, { stream: true });
        }

        try {
            const bodyJson: RawResponseData = JSON.parse(responseBody);
            const responseData: ResponseData = {
                success: bodyJson.success,
                result: bodyJson.result,
                error: ErrorCode[bodyJson.error]
            }

            if (response.ok) reslove(responseData);
            else reject(ErrorCode[bodyJson.error]);
        }
        catch (err) {
            reject(err);
        }
    });
}

function GetMe(): Promise<ResponseData> {
    return Get(Path.USER)
}

function GetInfo(appId: string = ""): Promise<ResponseData> {
    return Get(Path.INFO + appId);
}

function SetInfo(name: string, callback: string, description: string, permissions: number, id: string): Promise<ResponseData> {
    return Post(Path.INFO, { name, callback, description, permissions, id });
}

function AcceptApp(appId: string): Promise<ResponseData> {
    return Get(Path.APP + appId);
}

function DeleteApp(appId: string): Promise<ResponseData> {
    return Request("DELETE", Path.INFO + appId)
}

export default {
    GetMe,
    GetInfo,
    SetInfo,
    AcceptApp,
    DeleteApp
}