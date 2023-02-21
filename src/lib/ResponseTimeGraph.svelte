<script lang="ts">
  import {onDestroy, onMount} from "svelte/internal";
  import {fetchJsonApi} from "~/utils/api";

  export let dataUrl: string;
  export let width: number;

  let height: number = 150;
  $: plotWidth = width - 61;
  $: plotHeight = height - 51;

  let timeInterval;
  let rawData: Item[];
  let maxTime: [number, number, string] = [1, 1, "s"];
  let timeLabels: [number, string][];
  let graphPath: string;

  const totalSeconds = 86400;

  onMount(() => {
    recalcualteTime();
    timeInterval = setInterval(recalcualteTime, 10000);
  });

  onDestroy(() => {
    clearInterval(timeInterval);
  });

  async function recalcualteTime() {
    let data = (await fetchJsonApi(dataUrl, "GET")) as Item[];
    if (data.length == 0) return;
    rawData = data.map(x => {
      x.y /= 1e9;
      return x;
    });
    let m = Math.max(...rawData.map(x => x.y));
    let n = 1e9;
    let units = "s";
    if (~~m == 0) {
      m *= 1000;
      n = 1e6;
      units = "ms";
    }
    if (~~m == 0) {
      m *= 1000;
      n = 1e3;
      units = "µs";
    }
    if (~~m == 0) {
      m *= 1000;
      n = 1;
      units = "ns";
    }
    maxTime = [Math.ceil(m * 2) / 2, n, units];

    let d = new Date();
    timeLabels = calculateTimeLabels(d);
    graphPath = renderRawGraph(d);
  }

  function renderRawGraph(d: Date) {
    let d2 = new Date(d);
    d2.setHours(9);
    d2.setMinutes(0);
    d2.setSeconds(0);
    let d3 = new Date(d2);
    d3.setHours(10);
    return `M ${calculateXFromTime(d2, d)} 5 L ${calculateXFromTime(d3, d)} 10`;
  }

  function round1dp(d: number) {
    return Math.floor(d * 10) / 10;
  }

  function getRawSeconds(d: Date) {
    return Math.floor(d.getTime() / 1000);
  }

  function calculateXFromTime(t: Date, d: Date) {
    return (plotWidth / totalSeconds) * (totalSeconds - (getRawSeconds(d) - getRawSeconds(t)));
  }

  function calculateTimeLabels(d: Date): [number, string][] {
    let o = d.getHours() - (d.getHours() % 3);

    return [7, 6, 5, 4, 3, 2, 1, 0].map(x => {
      let d2 = new Date(d);
      d2.setHours(o - 3 * x);
      d2.setMinutes(0);
      d2.setSeconds(0);
      return [calculateXFromTime(d2, d) + 51, renderHours(d2.getHours())];
    });
  }

  function renderHours(h: number) {
    let b = false;
    if (h >= 12) {
      b = true;
      h -= 12;
    }
    if (h == 0) {
      h = 12;
    }
    return `${h < 10 ? "0" : ""}${h}:00${b ? "pm" : "am"}`;
  }

  interface Item {
    x: number;
    y: number;
  }
</script>

<div>
  {#if rawData}
    <p class="header">Response times</p>
    <div dir="ltr" class="chart">
      <svg version="1.1" xmlns="http://www.w3.org/2000/svg" {width} {height} viewBox="0 0 {width} {height}">
        <defs>
          <clipPath id="chart-oot8ea6-1-">
            <rect x="0" y="0" width={plotWidth} height={plotHeight} fill="none" />
          </clipPath>
        </defs>
        <g class="chart-axis chart-xaxis">
          {#each timeLabels as timeLabel}
            <path fill="none" class="chart-tick" d="M {timeLabel[0]} 112 L {timeLabel[0]} 122" opacity="1" />
          {/each}
          <path fill="none" class="chart-axis-line" stroke="#191C24" stroke-width="1" d="M 51 112.5 L {plotWidth + 51} 112.5" />
          <path stroke-dasharray="3,3" d="M 51 10 L {plotWidth + 51} 10" stroke-linecap="round" />
          <path stroke-dasharray="3,3" d="M 51 44 L {plotWidth + 51} 44" stroke-linecap="round" />
          <path stroke-dasharray="3,3" d="M 51 78 L {plotWidth + 51} 78" stroke-linecap="round" />
        </g>
        <g class="chart-series-group">
          <g
            class="chart-series chart-series-0 chart-areaspline-series"
            opacity="1"
            transform="translate(51,10) scale(1 1)"
            clip-path="url(#chart-oot8ea6-1-)"
          >
            {#if rawData}
              {JSON.stringify(rawData)}
            {/if}
            <path fill="none" d={graphPath} class="chart-graph" stroke="#70778C" stroke-width="1.5" stroke-linejoin="round" stroke-linecap="round" />
          </g>
          <g class="chart-axis-labels chart-xaxis-labels">
            {#each timeLabels as timeLabel}
              <text
                x={timeLabel[0]}
                style="color:#70778C;cursor:default;font-size:11px;font:12px Inter;letter-spacing:-0.2px;fill:#70778C;"
                text-anchor="middle"
                transform="translate(0,0)"
                y="136"
                opacity="1"
              >
                <tspan>{timeLabel[1]}</tspan>
              </text>
            {/each}
          </g>
          <g class="chart-axis-labels chart-yaxis-labels" data-z-index="7">
            <text
              x="46"
              style="color:#70778C;cursor:default;font-size:11px;font:12px Inter;letter-spacing:-0.2px;fill:#70778C;"
              text-anchor="end"
              transform="translate(0,0)"
              y="116"
              opacity="1"
            >
              <tspan>0{maxTime[2]}</tspan>
            </text>
            <text
              x="46"
              style="color:#70778C;cursor:default;font-size:11px;font:12px Inter;letter-spacing:-0.2px;fill:#70778C;"
              text-anchor="end"
              transform="translate(0,0)"
              y="82"
              opacity="1"
            >
              <tspan>{round1dp(maxTime[0] / 3)} {maxTime[2]}</tspan>
            </text>
            <text
              x="46"
              style="color:#70778C;cursor:default;font-size:11px;font:12px Inter;letter-spacing:-0.2px;fill:#70778C;"
              text-anchor="end"
              transform="translate(0,0)"
              y="48"
              opacity="1"
            >
              <tspan>{round1dp((maxTime[0] * 2) / 3)} {maxTime[2]}</tspan>
            </text>
            <text
              x="46"
              style="color:#70778C;cursor:default;font-size:11px;font:12px Inter;letter-spacing:-0.2px;fill:#70778C;"
              text-anchor="end"
              transform="translate(0,0)"
              y="14"
              opacity="1"
            >
              <tspan>{maxTime[0]} {maxTime[2]}</tspan>
            </text>
          </g>
          <g
            class="chart-label chart-tooltip chart-color-undefined"
            style="white-space:nowrap;font:12px Inter;pointer-events:none;"
            data-z-index="8"
            transform="translate(441,-9999)"
            opacity="0"
            visibility="hidden"
          >
            <path
              fill="#21242D"
              class="chart-label-box chart-tooltip-box"
              d="M 12.5 0.5 L 173.5 0.5 C 185.5 0.5 185.5 0.5 185.5 12.5 L 185.5 35.5 C 185.5 47.5 185.5 47.5 173.5 47.5 L 98.5 47.5 L 92.5 53.5 L 86.5 47.5 L 12.5 47.5 C 0.5 47.5 0.5 47.5 0.5 35.5 L 0.5 12.5 C 0.5 0.5 0.5 0.5 12.5 0.5"
              stroke="#2D313C"
              stroke-width="1"
            />
            <text x="8" data-z-index="1" y="20" style="color:#8A91A5;cursor:default;font-size:12px;fill:#8A91A5;">
              <tspan>Feb 19, 2023 at 09:38pm PST</tspan>
              <tspan style="font-size: 13px; font-weight: 500; fill: white;" x="8" dy="16">Response time: 0.29 s</tspan>
            </text>
          </g>
        </g>
      </svg>
    </div>
  {/if}
</div>

<style lang="scss">
  svg {
    font-family: "Lucida Grande", "Lucida Sans Unicode", Arial, Helvetica, sans-serif;
    font-size: 12px;
    position: static;

    .chart-axis.chart-xaxis path {
      stroke: #c0c0c0;
      stroke-width: 1;
    }
  }
</style>
