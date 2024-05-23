<script lang="ts">
    import { Routes, route, router } from "./Router";
    import { login } from "./Status";

    $: {
        switch (true) {
            case $route.INFO: {
                if ($login) break;
                router.Set(Routes.LOGIN);
            }
        }
    }
</script>
