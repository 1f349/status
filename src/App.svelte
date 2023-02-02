<script lang="ts">
  import { Router, Route, link } from "svelte-navigator";
  import LazyComponent from "./lib/LazyComponent.svelte";
  import { user } from "~/stores/user";
  import { popupCenterScreen } from "~/utils/window";
  import { isObject } from "~/utils/utils";
  import { API_URL } from "~/utils/api";

  function handleLoginClick() {
    popupCenterScreen(`${API_URL}/login?in_popup`, "Login", 800, 800, false);
  }

  function handleLogoutClick() {
    fetch(`${API_URL}/logout`, { method: "POST" }).then(function (resp) {
      if (resp.ok) {
        user.set(null);
      }
    });
  }

  function check_user() {
    const f = document.createElement("iframe");
    f.src = `${API_URL}/check`;
    f.style.display = "none";
    document.body.appendChild(f);
  }

  window.onmessage = function (event) {
    if (event.origin !== API_URL) return;
    if (isObject(event.data)) {
      console.log(event.data);
      if (isObject(event.data.user)) {
        let d = Object.assign({ sub: null, login: null, name: null, picture: null, admin: false }, event.data.user);
        if (d.sub === null || d.login === null || d.name === null || d.picture === null) {
          alert("Failed to log user in: the login data is structured correctly but probably corrupted");
          return;
        }
        user.set(d);
        return;
      }
    }
    alert("Failed to log user in: the login data was probably corrupted");
  };
</script>

<svelte:head>
  <title>{import.meta.env.VITE_TITLE}</title>
  <link rel="stylesheet" href="{import.meta.env.VITE_CSS_VAR}.light.css" media="screen" />
  <link rel="stylesheet" href="{import.meta.env.VITE_CSS_VAR}.dark.css" media="screen and (prefers-color-scheme: dark)" />
</svelte:head>

<div id="app-router">
  <Router primary={false}>
    <header>
      <div class="central-header">
        <a href="/" use:link>
          <h1>{import.meta.env.VITE_TITLE}</h1>
        </a>
        <nav>
          {#if $user}
            <div style="height:50px;">
              <a href="/admin/dashboard" use:link>Dashboard</a>
              <button on:click={handleLogoutClick}>Logout</button>
            </div>
          {:else}
            <button on:click={handleLoginClick}>Login</button>
          {/if}
        </nav>
      </div>
    </header>
    <Route path="admin">
      <LazyComponent component={() => import("./routes/Admin.svelte")} delayMs={500}>Loading...</LazyComponent>
    </Route>
    <Route path="service/*">
      <LazyComponent component={() => import("./routes/Service.svelte")} delayMs={500}>Loading...</LazyComponent>
    </Route>
    <Route path="/">
      <LazyComponent component={() => import("./routes/Home.svelte")} delayMs={500}>Loading...</LazyComponent>
    </Route>
  </Router>
  <footer>
    Powered by <a href="https://code.mrmelon54.com/melon/status" target="_blank" rel="noreferrer noopener nofollow">Status</a>
  </footer>
</div>

<style lang="scss">
  #app-router {
    max-width: min(100%, 1000px);
    margin: auto;

    > header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      height: 50px;
      min-height: 50px;
      max-height: 50px;
      background-color: var(--primary-main);
      border-radius: 0 0 var(--large-curve) var(--large-curve);

      > .central-header {
        padding: 0 32px;
        height: 50px;
        min-height: 50px;
        max-height: 50px;
        align-items: center;
        justify-content: space-between;
        display: flex;
        width: min(100%, 1000px);
        max-width: min(100%, 1000px);
        margin: auto;

        > .mobile-nav-toggle {
          display: none;
        }

        > nav {
          height: 100%;
          display: flex;
          align-items: center;
        }

        a {
          color: #eeeeee;

          &:hover {
            color: #ffffff;
          }

          > h1 {
            font-size: 26px;
            line-height: 50px;
          }
        }

        > nav > :global(a) {
          padding: 8px;
        }
      }
    }

    > footer {
      padding: 16px;
    }

    @media screen and (max-width: 600px) {
      > header {
        border-radius: 0;
        > .central-header {
          > .mobile-nav-toggle {
            display: flex;
            padding: 0;
            width: 50px;
            aspect-ratio: 1/1;
            align-items: center;
            justify-content: center;

            &.mobile-active {
              background-color: var(--primary-hover);
            }
          }

          > nav {
            position: fixed;
            top: 50px;
            left: 0;
            right: 0;
            height: auto;
            z-index: 9998;
            display: none;

            &.mobile-open {
              display: block;
              background-color: var(--bg-panel-action);
            }
          }

          > .mobile-shade.mobile-open {
            position: fixed;
            top: 50px;
            left: 0;
            right: 0;
            bottom: 0;
            background-color: rgba(0, 0, 0, 0.5);
            content: "";
            z-index: 9997;
          }
        }
      }
    }
  }
</style>
