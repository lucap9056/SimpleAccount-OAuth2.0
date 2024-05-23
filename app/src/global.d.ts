/// <reference types="svelte" />

declare type ResponseData = {
    success: boolean
    result?: string
    error?: string
}

declare type Route = {
    [key: string]: boolean
}

declare type User = {
    id?: number
    name?: string
    email?: string
    lastEditTime?: number
    registerTime?: number
    deletedTime?: number
}

declare type AppInfo = {
    id?: string
    name?: string
    client?: number
    secret?: string
    callback?: string
    description?: string
    permissions?: number
}

declare type AppPermissions = {
    permission_read_register_time?: boolean;
    permission_read_last_edit?: boolean;
    permission_read_email_address?: boolean;
    permission_read_username?: boolean;
}