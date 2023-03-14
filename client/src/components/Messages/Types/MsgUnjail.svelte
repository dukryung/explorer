<script>
    import {_} from 'svelte-i18n';
    import {Link} from 'svelte-routing';
    import {Token} from '../../../js/token';

    export let message;
</script>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('Unjail')}</h5></div>
                </div>
            </div>
            <div class="CardList">
                <table class="CardTable MsSend">
                    <thead>
                    </thead>
                    <tbody>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('validator_address')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark MSG">
                                <Link to="/validator/{message.validator_addr}">
                                    {message.validator_addr}
                                </Link>
                            </p>
                        </td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</section>