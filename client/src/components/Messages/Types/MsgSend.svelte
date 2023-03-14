<script>
    import {_} from 'svelte-i18n';
    import {Link} from 'svelte-routing';
    import {Token} from '../../../js/token';
    import Coin from '../../Cards/Coin.svelte';

    export let message;
</script>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('msg_send')}</h5></div>
                </div>
            </div>
            <div class="CardList">
                <table class="CardTable MsSend">
                    <thead>
                    </thead>
                    <tbody>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('msg_sender')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark MSG">
                                <Link to="/account/{message.from_address}">
                                    {message.from_address}
                                </Link>
                            </p>
                        </td>
                    </tr>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('msg_recipient')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark MSG">
                                <Link to="/account/{message.to_address}">
                                    {message.to_address}
                                </Link>
                            </p>
                        </td>
                    </tr>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('msg_amount')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            {#each message.amount as amount}
                                <p class="FS15 fwM Color_Dark">
                                    <!-- TODO get the actual token from the backend 
                                    {`${new Token({
                                        denom: amount.denom,
                                        amount: amount.amount,
                                        precision: network.networkToken.precision,
                                    }).string()}`}
                                -->
                                <Coin coin={amount}/>
                                </p>
                            {/each}
                        </td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</section>