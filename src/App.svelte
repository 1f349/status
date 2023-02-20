<script lang="ts">
  import {Router, Route, link} from "svelte-navigator";
  import LazyComponent from "./lib/LazyComponent.svelte";
  import NavItems from "./lib/NavItems.svelte";
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
        <NavItems />
      </div>
    </header>
    <Route path="maintenance">
      <LazyComponent component={() => import("./routes/Maintenance.svelte")} delayMs={500}>Loading...</LazyComponent>
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
    max-width: min(100%, 800px);
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
      }
    }

    > footer {
      padding: 16px;
    }

    @media screen and (max-width: 600px) {
      > header {
        border-radius: 0;
        > .central-header {
          > nav {
            position: fixed;
            top: 50px;
            left: 0;
            right: 0;
            height: auto;
            z-index: 9998;
            display: none;
          }
        }
      }
    }
  }
</style>
