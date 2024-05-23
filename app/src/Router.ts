import { writable, type Writable } from "svelte/store";
import API from "./API";
import { Alert, alertManager } from "./Alert/Struct";
import Loading from "./Loading/Main"
import Translations, { TranslationsGet } from "./Translations";
import Status from "./Status";
import loadings from "./Loading/Main";

let login: boolean = false;
Status.login.subscribe((l) => {
    login = l;
});

export class Hash {
    private dirs: string[];
    constructor(hash: string) {
        this.dirs = hash.replace(/^#/, "").split("/");
    }

    public Shift(): string {
        if (this.dirs.length > 0) {
            return this.dirs.shift();
        }
        return "";
    }
}

export const Routes: { [key: string]: string } = {
    LOGIN: "LOGIN",
    INDEX: "INDEX",
    INFO: "INFO",
    APP: "",
};

export const router = new class {
    public route: Writable<Route>;
    private hash: Hash;
    constructor() {
        const routes = this.CreateRoutes();
        this.route = writable(routes);

        this.hash = new Hash(location.hash);
        window.addEventListener("hashchange", this.HashChange.bind(this));
        this.HashChange();
    }

    public Set(...routes: string[]): void {
        if (routes.length < 1) return;
        let hash = "#" + routes[0].toLowerCase();

        for (let i = 1; i < routes.length; i++) {
            hash += "/" + routes;
        }
        location.hash = hash;
    }

    private CreateRoutes(): Route {
        const routes: Route = {};
        for (const route of Object.keys(Routes)) {
            routes[route] = false;
        }
        return routes;
    }

    private async HashChange() {
        const { CreateRoutes, Set } = this;
        const hash = new Hash(location.hash);
        this.hash = hash;

        const routes = CreateRoutes();

        const head = hash.Shift().toUpperCase();

        switch (head) {
            case Routes.LOGIN:
                location.replace("/#login");
                return;
            case Routes.INDEX:
                location.replace("/#index");
                return;
        }

        let exist = false;
        for (const key of Object.keys(Routes)) {
            if (head == Routes[key]) {
                exist = true;
                routes[key] = true;
            }
        }

        if (!exist) {
            if (login) {
                Set(Routes.INDEX);
            } else {
                Set(Routes.LOGIN);
            }
            return;
        }

        this.route.update(() => routes);
    }

    public Hash() {
        return this.hash;
    }
}

export const route = router.route;

export default {
    router,
    route,
    Routes
}