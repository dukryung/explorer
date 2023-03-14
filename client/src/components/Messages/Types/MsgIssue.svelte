<script>
    import {_} from 'svelte-i18n';
    import {Token} from '../../../js/token';


    export let message;
</script>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('msg_issue')}</h5></div>
                </div>
            </div>
            <div class="CardList">
                <table class="CardTable MsSend">
                    <thead>
                    </thead>
                    <tbody>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('msg_issuer')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark MSG">{message.owner_address}</p>
                        </td>
                    </tr>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('msg_symbol')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark MSG">{message.symbol }</p>
                        </td>
                    </tr>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('msg_description')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark MSG">{message.description }</p>
                        </td>
                    </tr>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('msg_precision')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark MSG">{message.precision }</p>
                        </td>
                    </tr>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('msg_coin')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                                <p class="FS15 fwM Color_Dark">
                                    {new Token({
                                        denom: message.coin.denom,
                                        amount: message.coin.amount,
                                        precision: message.precision})
                                        .string()
                                    }
                                </p>
                        </td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</section>