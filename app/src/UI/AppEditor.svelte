<script lang="ts">
    import Translations from "../Translations";
    import PERM from "../PERM";
    import Textbox from "../Components/Textbox.svelte";
    import ValidateInput from "../ValidateInput";
    import API from "../API";
    import loadings from "../Loading/Main";
    import { Alert, alertManager } from "../Alert/Struct";

    export let info: AppInfo;
    let permissions: AppPermissions = {};

    if (info.id) {
        permissions = PERM.Parse(info.permissions);
    }

    function AllPermissions(app: AppInfo): { key: string; value: boolean }[] {
        const perms = PERM.Parse(app.permissions);
        return Object.keys(perms).map((key) => ({ key, value: perms[key] }));
    }

    function handleSubmit() {
        if (info.name == "" || info.callback == "") return;

        const laoding = loadings.Append();
        API.SetInfo(
            info.name,
            info.callback,
            info.description,
            info.permissions,
            info.id,
        )
            .then((res) => {
                laoding.Remove();
                const result = JSON.parse(res.result);
                info.id = result.id;
                info.secret = result.secret;
            })
            .catch((err) => {
                laoding.Remove();
                alertManager.Add($Translations[err], Alert.Type.Error);
            });
    }

    function Cancel() {
        info = null;
    }

    function PermissionClick() {
        const permission: HTMLLIElement = this;
        const name = permission.dataset.name;
        const value = permission.dataset.value != "true";
        const checkbox: HTMLInputElement =
            permission.getElementsByTagName("input")[0];

        permission.dataset.value = value.toString();
        checkbox.checked = !value;
        permissions[name] = value;

        info.permissions = PERM.ToNumber(permissions);
    }

    function Copy() {
        this.select();
        document.execCommand("copy");
        alertManager.Add($Translations.edit_value_copy, Alert.Type.Normal);
    }

    function Reload() {
        location.reload();
    }
</script>

<div class="edit container">
    {#if info.id && info.secret}
        <div class="secret_label">{$Translations.edit_label_secret}</div>
        <form class="secret_form" on:reset={Reload}>
            <div class="secret-group">
                <label for="app_id">Token</label>
                <input id="app_id" value={info.id} readonly on:click={Copy} />
                <div class="hint"></div>
            </div>
            <div class="secret-group">
                <label for="app_secret">Secret</label>
                <textarea id="app_secret" on:click={Copy} value={info.secret}
                ></textarea>
                <div class="hint"></div>
            </div>
            <button type="reset">{$Translations.confirm}</button>
        </form>
    {:else}
        <h1>{info.id ? $Translations.edit_app : $Translations.create_app}</h1>
        <form on:submit|preventDefault={handleSubmit} on:reset={Cancel}>
            <Textbox
                name="name"
                bind:input_value={info.name}
                hint={$Translations.edit_app_name_limit}
                validate={ValidateInput.AppName}
            />
            <Textbox
                name="callback"
                bind:input_value={info.callback}
                hint={$Translations.edit_callback_hint}
                validate={ValidateInput.AppCallback}
                placeholder={$Translations.edit_callback_placeholder}
            />
            <textarea class="app_description" bind:value={info.description} />
            <div class="form-group permissions">
                <ul>
                    {#each AllPermissions(info) as { key, value }}
                        <!-- svelte-ignore a11y-click-events-have-key-events -->
                        <li
                            class="permission"
                            data-name={key}
                            data-value={value}
                            on:click={PermissionClick}
                        >
                            <div class="permission_name">
                                {$Translations[key]}
                            </div>
                            <input type="checkbox" checked={value} />
                            <div class="checkbox-custom"></div>
                        </li>
                    {/each}
                </ul>
            </div>
            <div class="form-group buttons">
                <button type="submit">
                    {info.id ? $Translations.confirm : $Translations.create}
                </button>
                <button type="reset">{$Translations.cancel}</button>
            </div>
        </form>
    {/if}
</div>

<style>
    .edit.container {
        max-width: 500px;
        margin: 0 auto;
        padding: 20px;
        background-color: #f7f7f7;
        border-radius: 10px;
        box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
    }

    .app_description {
        width: 100%;
        resize: none;
        height: 90px;
        padding: 10px;
        border-radius: 5px;
        border: 1px solid #ccc;
        box-sizing: border-box;
    }

    .form-group.permissions {
        width: 50%;
        display: inline-block;
        text-align: right;
    }

    .form-group.buttons {
        float: right;
        width: 50%;
        display: grid;
        justify-content: center;
    }

    .permission {
        margin-bottom: 10px;
        cursor: pointer;
    }

    .permission_name {
        display: inline-block;
        vertical-align: middle;
        margin-right: 10px;
        color: #333;
    }

    .permissions ul {
        padding-left: 0;
        list-style-type: none;
    }

    .buttons button[type="submit"] {
        background-color: #4caf50;
        color: white;
    }

    .buttons button {
        border: none;
        padding: 10px 20px;
        border-radius: 5px;
        cursor: pointer;
        width: 140px;
        margin-bottom: 10px;
        display: grid;
        background-color: #6c757d;
        margin-top: 20px;
        color: white;
    }

    .buttons button:hover {
        background-color: #5a6268;
    }

    .buttons button[type="submit"]:hover {
        background-color: #45a049;
    }

    input[type="checkbox"] {
        display: none;
    }

    .checkbox-custom {
        display: inline-block;
        width: 20px;
        height: 20px;
        background-color: #eee;
        border-radius: 5px;
        vertical-align: middle;
        position: relative;
    }

    .checkbox-custom::after {
        content: "";
        width: 6px;
        height: 12px;
        border: solid #333;
        border-width: 0 2px 2px 0;
        transform: rotate(45deg);
        position: absolute;
        top: 3px;
        left: 7px;
        opacity: 0;
    }

    input[type="checkbox"]:checked + .checkbox-custom::after {
        opacity: 1;
    }

    .secret_form {
        width: 450px;
        margin: auto;
    }

    .secret_label {
        line-height: 32px;
        font-weight: bold;
    }

    .secret-group {
        margin-bottom: 20px;
    }

    .secret-group label {
        width: 120px;
        text-align: right;
        float: left;
        line-height: 36px;
    }

    #app_id {
        width: 300px;
        height: 36px;
        font-size: 16px;
        padding: 2px;
    }

    #app_secret {
        resize: none;
        user-select: all;
        width: 300px;
        height: 36px;
        line-height: 16px;
        font-size: 16px;
        padding: 2px;
    }

    .secret_form button {
        border: none;
        padding: 10px 20px;
        border-radius: 5px;
        cursor: pointer;
        width: 140px;
        margin-bottom: 10px;
        background-color: #4caf50;
        color: white;
    }

    .secret_form button[type="reset"]:hover {
        background-color: #45a049;
    }
</style>
