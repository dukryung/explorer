<script>
    import { onMount } from 'svelte';
    import { Token } from 'js/token';

    export let coin;

    let onLoad = false;
    let token;

    onMount(() => {
        // TODO Use cache to keep coins info
        klaatoo.singleRequest({
            method: 'token.bydenom',
            params: [coin.denom],
            id: klaatoo.generateRequestId(),
            success: data => {
                if (data.symbol) {
                    token = new Token({
                        symbol: data.symbol,
                        precision: data.precision,
                        amount: coin.amount
                    });
                } else {
                    token = new Token({
                        symbol: coin.denom,
                        amount: coin.amount,
                        precision: 0
                    });
                }
                onLoad = true;
            },
            error: error => {
                console.error(error);
            }
        });
    });
</script>

<span>{onLoad ? token.string() : ''}</span>
