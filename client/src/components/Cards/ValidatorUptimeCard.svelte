<script>
    import {onDestroy, onMount} from 'svelte';
    import {_} from 'svelte-i18n';
    import {Link, navigate} from 'svelte-routing';
    import LoadingIndicator from 'Components/Labels/LoadingIndicator';
    import Failed from '../../assets/Failed.svelte';
    import Success from '../../assets/Success.svelte';

    export let address;

    let uptime;
    let onLoad = false;
    const requestId = klaatoo.generateRequestId();

    onMount(() => {
      klaatoo.subscribe({
        method: 'validator.uptime',
        params: [
          address,
        ],
        id: requestId,
        success: (data) => {
          if (data != null) {
            uptime = data;
          }
          onLoad = true;
        },
        error: (error) => {
          console.error(error);
          navigate('/error', {replace: false});
        },
      });
    });

    onDestroy(() => {
      klaatoo.unsubscribe({
        id: requestId,
        method: 'validator.uptime',
        params: [],
      });
    });
</script>

<style>
    .grid {
        width: auto;
        margin-top:10px;
        margin-bottom: 10px;
    }
</style>


<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('uptime')}</h5></div>
                    <h6 class="TW FS13"><span class="TW FS13">Latest 100 blocks</span></h6>
                </div>
            </div>
            {#if onLoad}
                <div class="CardList">
                    <table class="CardTable AcInfo">
                        <colgroup>
                            <col class="LeftCol">
                            <col class="RightCol">
                        </colgroup>
                        <thead>
                        </thead>
                        <tbody>
                        </tbody>
                    </table>
                    <div class="grid">
                        {#each Array(100) as _, i}
                            {#if uptime.blocks.find((block) => block.height === uptime.latest_height - i - 1)}
                                <Link to="/block/{uptime.latest_height-i-1}">
                                    <Failed/>
                                </Link>
                            {:else}
                                <Link to="/block/{uptime.latest_height - i -1}">
                                    <Success/>
                                </Link>
                            {/if}
                        {/each}
                    </div>
                </div>
            {:else}
                <LoadingIndicator/>
            {/if}
        </div>
    </div>
</section>


