<script>
    import {onDestroy} from 'svelte';
    import {_} from 'svelte-i18n';
    import {stats} from 'js';

    let validatorTotal;

    const unsubscribe = stats.subscribe((v) => validatorTotal = v.validator_total);
    onDestroy(unsubscribe);
</script>

<div class="card">
    <div class="CardBackground"></div>
    <div class="CardContent">
        <p class="CardCategory">
            {$_('total_validator')}
        </p>
        <img src="images/b_vail.png">
        <h3 class="CardHeading">
            {validatorTotal}
        </h3>
    </div>
</div>
