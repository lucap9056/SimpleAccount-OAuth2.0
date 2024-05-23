
export default class {
    public static Parse(permNum: number): AppPermissions {
        const permStr = permNum.toString(2).padStart(4, "0");
        return {
            permission_read_register_time: permStr[0] == "1",
            permission_read_last_edit: permStr[1] == "1",
            permission_read_email_address: permStr[2] == "1",
            permission_read_username: permStr[3] == "1",
        }
    }

    public static ToNumber(perms: AppPermissions): number {
        const Permissions = [
            perms.permission_read_username,
            perms.permission_read_email_address,
            perms.permission_read_last_edit,
            perms.permission_read_register_time
        ];
        let permNum = 0
        for (let i = 0; i < Permissions.length; i++) {
            const perm = Permissions[i];
            if (perm) {
                permNum += Math.pow(2, i);
            }
        }

        return permNum;
    }
}