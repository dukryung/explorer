<script>
    import {onDestroy} from 'svelte';
    import {_} from 'svelte-i18n';
    import {stats} from 'js';

    let blockBestTime;

    const unsubscribe = stats.subscribe((v) => blockBestTime = v.block_min_time);
    onDestroy(unsubscribe);
</script>



<div class="card">
    <div class="CardBackground"></div>
    <div class="CardContent">
        <p class="CardCategory">
            {$_('block_best_time')}
        </p>
        <img src="/images/b_time.png">
        <h3 class="CardHeading">
            {blockBestTime / 1000}{$_('second')}
        </h3>
    </div>
</div>