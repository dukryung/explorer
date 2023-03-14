<script>
    import {Token} from 'js/token';

    export let token;

    $: if (token) {
      if (token.coin.denom === network.networkToken.base) {
        token.symbol = network.networkToken.symbol;
        token.precision = network.networkToken.precision;
        token.description = network.networkToken.description;
        token.amount = token.coin.amount;
        token.src = '/images/networkTokenIcon.png';
      } else {
        token.symbol = token.symbol;
        token.description = token.description;
        token.precision = token.precision;
        token.amount = token.coin.amount;
        token.src = '/images/noIcon.png';
      }
    }
</script>


<tr>
    <td class="px-2 py-1 TextAL">
        <div class="TokenIcon">
            <img src="{token.src}">
        </div>
        <p class="FS18 pT10 pL65">
            {token.symbol}
        </p>
        <!-- Tx Sender -->
        <p class="FS12 Color_Gray pL65 TxAddress">
            {token.description}
        </p>
    </td>
    <td class="px-2 py-1 TextAR TokenMo">
        <p class="FS15 fwB Color_Dark">
            <!--{new Token(token.symbol, token.amount, token.precision).string()}-->
            {`${new Token({
                symbol: token.symbol,
                amount: token.amount,
                precision: token.precision,
            }).string()}`}
        </p>
    </td>
</tr>