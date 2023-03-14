<script>
    import {_} from 'svelte-i18n';
    import {afterUpdate, onDestroy, onMount} from 'svelte';
    import {Chart, registerables} from 'chart.js/dist/chart.esm';

    Chart.register(...registerables);
    let chart = null;
    let chartRef;

    // chart
    const height = 170;
    const backgroundColor = 'rgba(13,111,229,0.1)';
    const borderColor = 'rgba(13,111,229,1)';
    const borderWidth = 2;

    // date
    let firstDate;
    let lastDate;

    $:if (datasets && options && chartRef && chart !== null) {
      chart.data = datasets;
      chart.type = 'line';
      chart.options = options;
      chart.update();
    }

    afterUpdate(() => {
      if (!chart) return;
      chart.data = datasets;
      chart.type = 'line';
      chart.options = options;
      chart.update();
    });

    const requestId = klaatoo.generateRequestId();
    onMount(() => {
      chart = new Chart(chartRef, {
        type: 'line',
        data: datasets,
        options: options,
      });

      klaatoo.subscribe({
        method: 'stats.tx',
        params: [],
        id: requestId,
        success: (data) => {
          const l = [];
          const d = [];
          for (const stats of data) {
            l.push(new Date(stats.time_stamp).toDateString().substr(4, 6));
            d.push(stats.tx_count);
          }
          firstDate = l[0];
          lastDate = l[l.length - 1];

          datasets = getData(l, d);
        },
        error: (error) => {
          console.error(error);
        },
      });
    });
    onDestroy(() => {
      klaatoo.unsubscribe({
        method: 'stats.tx',
        id: requestId,
      });
    });

    function getData(labels, datasets) {
      return {
        labels: labels,
        datasets: [{
          fill: true,
          pointRadius: 1,
          lineTension: 0.1,
          data: datasets,
          backgroundColor: backgroundColor,
          borderColor: borderColor,
          borderWidth: borderWidth,
        }],
      };
    }

    let datasets = {};
    const options = {
      animation: false,
      maintainAspectRatio: false,
      scales: {
        x: {
          display: true,
        },
        y: {
          beginAtZero: false,
          display: true,
        },
      },
      plugins: {
        legend: {
          display: false,
        },
        tooltip: {
          mode: 'nearest',
          intersect: false,
          callbacks: {
            label: function(context) {
              let label = context.dataset.label || '';
              if (label) {
                label += ': ';
              }
              if (context.parsed.y !== null) {
                label += context.parsed.y + ' txs';
              }
              return label;
            },
          },
        },
      },
    };
</script>
<section class="BoxSc">
    <div class="BoxScWrap">
        <div class="titleTxt">
            <h5>
                {$_('latest_txs_14days')}
            </h5>
            <div class="tLine"></div>
        </div>
        <!-- Chart -->
        <div class="chart-container">
            <div class="ChartContainerWarp">
                <canvas id="myChart" bind:this={chartRef}></canvas>
            </div>
        </div>
        <!-- Chart date -->
        <!-- <div>
             {#if firstDate}
                 <p class="FL FS12 pD20 pL30 FwM MT-10">
                     {firstDate}
                 </p>
             {/if}
             {#if lastDate}
                 <p class="FR FS12 pD20 pR30 FwM MT-10">
                     {lastDate}
                 </p>
             {/if}
         </div>-->
    </div>
</section>
