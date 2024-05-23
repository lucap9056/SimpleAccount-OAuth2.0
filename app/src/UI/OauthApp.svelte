<script lang="ts">
    import API from "../API";
    import { Alert, alertManager } from "../Alert/Struct";
    import loadings from "../Loading/Main";
    import Translations from "../Translations";
    import PERM from "../PERM";

    let appInfo: AppInfo = null;
    let loaded = false;

    const query = new URLSearchParams(location.search);
    const appId = query.get("app");

    (async () => {
        const loading = loadings.Append();
        if (appId) {
            try {
                const res = await API.GetInfo(appId);
                if (res.success) {
                    appInfo = JSON.parse(res.result);
                } else {
                    alertManager.Add(
                        $Translations[res.error],
                        Alert.Type.Error,
                    );
                }
            } catch (err) {
                alertManager.Add($Translations[err], Alert.Type.Error);
            }
        }
        loading.Remove();
        loaded = true;
    })();

    function EnabledPermissions(app: AppInfo): string[] {
        console.log(app);
        const permissions = PERM.Parse(app.permissions || 0);
        return Object.keys(permissions).filter((key) => permissions[key]);
    }

    function Cancel() {
        window.close();
    }

    function Confirm() {
        const loading = loadings.Append();
        API.AcceptApp(appId)
            .then((res) => {
                if (res.success) {
                    location.replace(res.result);
                }
            })
            .catch((err) => {
                alertManager.Add($Translations[err], Alert.Type.Error);
                loading.Remove();
            });
    }
</script>

{#if loaded}
    {#if appInfo === null}
        <div class="null_app">App does not exist.</div>
    {:else}
        <form
            class="app_info"
            on:submit|preventDefault={Confirm}
            on:reset={Cancel}
        >
            <div class="app_name">{appInfo.name}</div>
            <div class="app_description">{appInfo.description}</div>
            <div class="app_permissions">
                <ul>
                    <p>{$Translations.permissions_requested}</p>
                    <li class="permission">
                        {$Translations.permission_read_id}
                    </li>
                    {#each EnabledPermissions(appInfo) as name}
                        <li class="permission">{$Translations[name]}</li>
                    {/each}
                </ul>
            </div>

            <div class="options">
                <button type="reset">{$Translations.cancel}</button>
                <button type="submit">{$Translations.confirm}</button>
            </div>
        </form>
    {/if}
{/if}

<style>
    .null_app {
        padding: 20px;
        margin: 20px;
        background-color: #ffdddd;
        border: 1px solid #ff0000;
        color: #d8000c;
        text-align: center;
        font-size: 18px;
    }

    .app_info {
        padding: 20px;
        margin: 20px;
        background-color: #f9f9f9;
        border: 1px solid #ddd;
        border-radius: 8px;
    }

    .app_name {
        font-size: 24px;
        font-weight: bold;
        margin-bottom: 10px;
    }

    .app_description {
        word-wrap: break-word;
    }

    .app_permissions p {
        margin-bottom: 0;
        text-align: left;
    }

    .app_permissions {
        margin: 10px 0;
    }

    .app_permissions ul {
        list-style-type: none;
        padding: 0;
    }

    .app_permissions .permission {
        padding: 8px;
        background-color: #e9e9e9;
        margin-bottom: 5px;
        border-radius: 4px;
    }

    .options {
        display: flex;
        justify-content: space-between;
        margin-top: 20px;
    }

    .options button {
        padding: 10px 20px;
        font-size: 16px;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }

    .options button[type="reset"] {
        background-color: #f44336;
        color: white;
    }

    .options button[type="submit"] {
        background-color: #4caf50;
        color: white;
    }
</style>
