<script lang="ts">
  import Check from "~/icons/Check.svelte";
  import Cross from "~/icons/Cross.svelte";
  import LazyDelay from "~/lib/LazyDelay.svelte";
  import {fetchStatus, fetchStatusOfService} from "~/utils/api";

  let fetching = fetchStatus() as Promise<Group[]>;
  fetching.then(x => {
    x.map(y => {
      y.services.map(z => {
        z.beans = fetchStatusOfService(z.id) as Promise<BeanStatus>;
        return z;
      });
      return y;
    });
  });

  interface Group {
    id: number;
    name: string;
    services: Service[];
  }

  interface Service {
    id: number;
    name: string;
    beans: Promise<BeanStatus>;
  }

  interface BeanStatus {
    current: Bean;
    beans: Bean[];
  }

  interface Bean {
    state: number;
    time: number;
  }
</script>

{#await fetching}
  <LazyDelay delayMs={500} />
{:then x}
  <div class="wrapper-block">
    {#each x as y (y.id)}
      <div class="group-item rounded-card">
        <h1>{y.name}</h1>
        {#each y.services as z (z.id)}
          <div class="service-item">
            <div class="title-row">
              {#await z.beans then zBeans}
                {#if zBeans.current.state == 1}
                  <div class="check-icon">
                    <Check size={18} />
                  </div>
                {:else if zBeans.current.state == 2}
                  <div class="cross-icon">
                    <Cross size={18} />
                  </div>
                {:else}{/if}
              {/await}
              <h2>{z.name}</h2>
            </div>
            <div class="beans">
              {#await z.beans}
                <LazyDelay delayMs={500} />
              {:then zBeans}
                {#each zBeans.beans as bean}
                  <div class="bean-view {bean.state == 1 ? 'bean-on' : 'bean-off'}" />
                {/each}
              {:catch}
                <div>Error loading</div>
              {/await}
            </div>
            <div class="graph" />
          </div>
        {/each}
      </div>
    {/each}
  </div>
{:catch}
  <div>Error loading</div>
{/await}

<style lang="scss">
  .wrapper-block {
    display: flex;
    gap: 12px;
    flex-direction: column;
    margin: 12px;

    .group-item {
      padding: 0.75rem;
      border-radius: 12px;
      box-shadow: 0 4px 7px 0 #00000008;
      border: 1px solid #21242d;

      h1 {
        margin: 0;
        line-height: 150%;
      }

      .service-item {
        display: flex;
        flex-direction: column;

        .title-row {
          display: flex;
          flex-direction: row;
          align-items: center;

          h2 {
            margin: 0;
            line-height: 150%;
          }
        }
      }
    }
  }

  .check-icon {
    background: var(--primary-main);
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: 50%;
  }

  .cross-icon {
    background: var(--secondary-main);
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: 50%;
  }
</style>
