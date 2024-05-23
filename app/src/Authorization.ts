const storage = window.localStorage;
delete window.localStorage;

const TOKEN = "tk";
const TEMP_TOKEN = "_tk";

function GetToken(): string {
    return storage.getItem(TEMP_TOKEN) || storage.getItem(TOKEN) || "";
}

function GetUser(): Promise<User> {
    return new Promise((resolve, reject) => {
        try {
            const token = storage.getItem(TOKEN).split(/\./g);
            const userStr = atob(token[1])
            const user: User = JSON.parse(userStr);
            resolve(user);
        }
        catch {
            reject();
            return;
        }
    });
}

export default {
    GetToken,
    GetUser
}