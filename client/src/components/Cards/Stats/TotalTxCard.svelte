<script>
    import {onDestroy} from 'svelte';
    import {_} from 'svelte-i18n';
    import {stats} from 'js';

    let totalTx;

    const unsubscribe = stats.subscribe((v) => totalTx = v.tx_total);
    onDestroy(unsubscribe);
</script>


<div class="card">
    <div class="CardBackground"></div>
    <div class="CardContent">
        <p class="CardCategory">
            {$_('total_transactions')}
        </p>
        <img src="/images/b_txs.png">
        <h3 class="CardHeading">
            {totalTx}
        </h3>
    </div>
</div>
