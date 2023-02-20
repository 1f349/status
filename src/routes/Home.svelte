<script lang="ts">
  import Check from "~/icons/Check.svelte";
  import Cross from "~/icons/Cross.svelte";
  import LazyDelay from "~/lib/LazyDelay.svelte";
  import {fetchStatus, fetchStatusOfService} from "~/utils/api";

  let innerWidth: number;
  $: days = calculateNumberOfDays(innerWidth);
  $: fetching = fetchStatus() as Promise<Group[]>;
  $: fetchBeans(days);
  let fetched = false;

  function calculateNumberOfDays(w: number) {
    if (w >= 600) return 90;
    if (w >= 400) return 60;
    return 30;
  }

  function fetchBeans(d: number) {
    fetching.then(x => {
      x.map(y => {
        y.services.map(z => {
          z.beans = fetchStatusOfService(z.id, d) as Promise<BeanStatus>;
          return z;
        });
        return y;
      });
      fetching = fetching;
      fetched = true;
    });
  }

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

<svelte:window bind:innerWidth />

<div class="wrapper-block">
  {#await fetching}
    <LazyDelay delayMs={500} />
  {:then x}
    {#each x as y (y.id)}
      <div class="group-item rounded-card">
        <h1>{y.name}</h1>
        {#each y.services as z (z.id)}
          <div class="service-item">
            <div class="title-row">
              {#await z.beans then zBeans}
                {#if zBeans.current.state == 1}
                  <div class="check-icon">
                    <Check size={14} />
                  </div>
                {:else if zBeans.current.state == 2}
                  <div class="cross-icon">
                    <Cross size={14} />
                  </div>
                {/if}
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
            <div class="bean-timeline">
              <div>{days} days ago</div>
              <div>Today</div>
            </div>
            <div class="graph" />
          </div>
        {/each}
      </div>
    {/each}
  {:catch}
    <div>Error loading</div>
  {/await}
</div>

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
        font-size: 22px;
      }
    }
  }

  .service-item {
    display: flex;
    flex-direction: column;
    margin-top: 1.5rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid #21242d;

    &:last-child {
      border-bottom: none;
    }

    .title-row {
      display: flex;
      flex-direction: row;
      align-items: center;

      h2 {
        margin: 0;
        line-height: 150%;
        font-size: 15px;
      }
    }

    .beans {
      display: grid;
      margin-top: 0.75rem;
      grid-template-columns: repeat(30, 1fr);
      gap: 1px;

      @media screen and (min-width: 400px) {
        & {
          grid-template-columns: repeat(60, 1fr);
        }
      }

      @media screen and (min-width: 600px) {
        & {
          grid-template-columns: repeat(90, 1fr);
        }
      }

      .bean-view {
        height: 32px;
        border-radius: 1px;

        &:first-child {
          border-radius: 4px 1px 1px 4px;
        }

        &:last-child {
          border-radius: 1px 4px 4px 1px;
        }

        &.bean-on {
          background-color: var(--bean-green);

          &:hover {
            background-color: var(--bean-green-hover);
          }
        }

        &.bean-off {
          background-color: var(--bean-red);

          &:hover {
            background-color: var(--bean-red-hover);
          }
        }
      }
    }

    .bean-timeline {
      display: flex;
      justify-content: space-between;
      margin-top: 0.75rem;
      margin-inline: 0.25rem;
      font-size: 13px;
      color: #8a91a5;
    }
  }

  .check-icon {
    background: var(--primary-main);
    display: flex;
    align-items: center;
    justify-content: center;
    width: 18px;
    height: 18px;
    margin: 3px;
    border-radius: 50%;
  }

  .cross-icon {
    background: var(--secondary-main);
    display: flex;
    align-items: center;
    justify-content: center;
    width: 18px;
    height: 18px;
    margin: 3px;
    border-radius: 50%;
  }
</style>
