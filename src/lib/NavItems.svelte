<script lang="ts">
  import {link, useLocation} from "svelte-navigator";

  const loc = useLocation();
</script>

<nav>
  <a href="/" class={$loc.pathname == "/" ? "nav-sel" : ""} use:link>Status</a>
  <a href="/maintenance" class={$loc.pathname == "/maintenance" ? "nav-sel" : ""} use:link>Maintenance</a>
</nav>

<style lang="scss">
  nav {
    height: 100%;
    display: flex;
    align-items: center;
    gap: 6px;

    a {
      color: #eeeeee;
      padding: 6px;
      border-radius: 6px;

      &:hover {
        color: #ffffff;
      }

      &.nav-sel {
        background-color: var(--primary-dark);
      }
    }
  }
</style>
