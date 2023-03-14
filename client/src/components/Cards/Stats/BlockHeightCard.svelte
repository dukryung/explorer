<script>
    import {onDestroy} from 'svelte';
    import {_} from 'svelte-i18n';
    import {stats} from 'js';

    let blockHeight;

    const unsubscribe = stats.subscribe((v) => blockHeight = v.block_height);
    onDestroy(unsubscribe);
</script>


<div class="card">
    <div class="CardBackground"></div>
    <div class="CardContent">
        <p class="CardCategory">
            {$_('block_height')}
        </p>
        <img src="/images/b_blocks.png">
        <h3 class="CardHeading">
            {blockHeight}
        </h3>
    </div>
</div>
