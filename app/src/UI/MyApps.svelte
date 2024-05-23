<script lang="ts">
    import API from "../API";
    import { Alert, alertManager } from "../Alert/Struct";
    import Authorization from "../Authorization";
    import loadings from "../Loading/Main";
    import AppEditor from "./AppEditor.svelte";
    import Translations from "../Translations";
    import PERM from "../PERM";

    let view: AppInfo = null;
    let edit: AppInfo = null;
    let infos: AppInfo[] = [];

    const loading = loadings.Append();

    (async () => {
        API.GetInfo()
            .then((res) => {
                if (res.success) {
                    infos = JSON.parse(res.result);
                } else {
                    alertManager.Add(res.error, Alert.Type.Error);
                }
                loading.Remove();
            })
            .catch((err) => {
                loading.Remove();
                alertManager.Add($Translations[err], Alert.Type.Error);
            });
    })();

    function View() {
        const i = this.dataset.index;
        view = infos[i];
    }

    function CancelView() {
        view = null;
    }

    function Edit() {
        edit = JSON.parse(JSON.stringify(view));
    }

    function NewApp() {
        edit = {
            name: "",
            callback: "",
            permissions: 0,
        };
    }

    let del = false;
    function DeleteCheck() {
        del = true;
    }

    function CancelDelete() {
        del = false;
    }

    function ConfirmDelete() {
        const loading = loadings.Append();
        API.DeleteApp(view.id)
            .then(() => {
                location.reload();
                loading.Remove();
            })
            .catch((err) => {
                alertManager.Add($Translations[err], Alert.Type.Error);
                loading.Remove();
            });
    }

    function EnabledPermissions(app: AppInfo): string[] {
        const permissions = PERM.Parse(app.permissions || 0);
        return Object.keys(permissions).filter((key) => permissions[key]);
    }
</script>

{#if edit !== null}
    <AppEditor bind:info={edit} />
{:else if view !== null}
    <div class="view container" style="pointer-events: {del ? 'none' : 'all'};">
        <div class="form-group">
            <div class="info_name">{view.name}</div>
        </div>
        <div class="form-group">
            <a target="_blank" href={view.callback}>
                <div class="info_callback">{view.callback}</div>
            </a>
        </div>
        <div class="description">{view.description}</div>
        <div class="permissions-group">
            <ul>
                <li class="permission">{$Translations.permission_read_id}</li>
                {#each EnabledPermissions(view) as name}
                    <li class="permission">
                        {$Translations[name]}
                    </li>
                {/each}
            </ul>
        </div>
        <div class="button-group">
            <button data-type="delete" on:click={DeleteCheck}>
                {$Translations.delete}
            </button>
            <button data-type="edit" on:click={Edit}>
                {$Translations.edit}
            </button>
            <button data-type="back" on:click={CancelView}>
                {$Translations.back}
            </button>
        </div>
    </div>
    {#if del}
        <div class="delete-confirmation">
            <h2>{$Translations.confirm_delete_title}</h2>
            <p>{$Translations.confirm_delete_message}</p>
            <div class="buttons">
                <button class="cancel-button" on:click={CancelDelete}>
                    {$Translations.cancel}</button
                >
                <button class="confirm-button" on:click={ConfirmDelete}>
                    {$Translations.confirm}
                </button>
            </div>
        </div>
    {/if}
{:else}
    <div class="list container">
        <h1>{$Translations.myapps}</h1>
        <ul>
            {#each infos as info, i}
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <li class="app" on:click={View} data-index={i}>
                    {info.name}
                </li>
            {/each}
        </ul>
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div class="create_app" on:click={NewApp}>
            {$Translations.myapps_new}
        </div>
    </div>
{/if}

<style>
    .container {
        width: 600px;
        margin: 0 auto;
        padding: 20px;
        background-color: #fff;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        border-radius: 8px;
    }

    .list.container {
        width: 400px;
    }

    h1 {
        font-size: 2em;
        color: #333;
        border-bottom: 2px solid #ddd;
        padding-bottom: 10px;
        margin-bottom: 20px;
    }

    .form-group {
        margin-bottom: 20px;
    }

    .info_name,
    .info_callback {
        font-size: 1.2em;
        padding: 10px;
        background-color: #f9f9f9;
        border: 1px solid #ccc;
        border-radius: 4px;
        margin-bottom: 10px;
    }

    a .info_callback {
        color: #007bff;
        text-decoration: none;
    }

    a .info_callback:hover {
        text-decoration: underline;
    }

    .permissions-group,
    .button-group {
        display: inline-block;
        vertical-align: top;
    }

    .permissions-group {
        width: 375px;
        margin-right: 20px;
    }

    .permissions-group ul {
        margin: 0;
    }

    .permissions-group li {
        padding: 10px;
        margin: 5px 0;
        background-color: #e7e7e7;
        border-radius: 4px;
    }

    .button-group {
        width: 200px;
        text-align: right;
    }

    .button-group button {
        padding: 10px 20px;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        margin: 5px 0;
        font-size: 1em;
        width: 100%;
    }

    .button-group button[data-type="delete"] {
        background-color: #dc3545;
        color: #fff;
        transition: background-color 0.3s ease;
    }

    .button-group button[data-type="delete"]:hover {
        background-color: #c82333;
    }

    button[data-type="back"] {
        background-color: #6c757d;
        color: #fff;
        transition: background-color 0.3s ease;
    }

    button[data-type="back"]:hover {
        background-color: #5a6268;
    }

    button[data-type="edit"] {
        background-color: hsl(211, 100%, 50%);
        color: #fff;
        transition: background-color 0.3s ease;
    }

    button[data-type="edit"]:hover {
        background-color: #0056b3;
    }

    ul {
        list-style-type: none;
        padding: 0;
    }

    li.app {
        padding: 10px;
        margin: 5px 0;
        background-color: #e7e7e7;
        cursor: pointer;
        border-radius: 4px;
        transition: background-color 0.3s ease;
    }

    li.app:hover {
        background-color: #d0d0d0;
    }

    .create_app {
        display: inline-block;
        padding: 10px 20px;
        background-color: #007bff;
        color: #fff;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        margin-top: 20px;
        transition: background-color 0.3s ease;
    }

    .create_app:hover {
        background-color: #0056b3;
    }

    .delete-confirmation {
        max-width: 400px;
        margin: 100px auto;
        padding: 20px;
        background-color: #fff;
        border: 1px solid #ddd;
        border-radius: 10px;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        text-align: center;
    }

    .delete-confirmation h2 {
        margin-top: 0;
        font-size: 24px;
        color: #d9534f;
    }

    .delete-confirmation p {
        font-size: 16px;
        color: #333;
        margin: 20px 0;
    }

    .delete-confirmation .buttons {
        display: flex;
        justify-content: space-between;
    }

    .delete-confirmation .buttons button {
        width: 48%;
        padding: 10px;
        font-size: 16px;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }

    .delete-confirmation .cancel-button {
        background-color: #ccc;
        color: #333;
    }

    .delete-confirmation .cancel-button:hover {
        background-color: #bbb;
    }

    .delete-confirmation .confirm-button {
        background-color: #d9534f;
        color: white;
    }

    .delete-confirmation .confirm-button:hover {
        background-color: #c9302c;
    }
</style>
