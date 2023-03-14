<script>
    import {_} from 'svelte-i18n';
    import {Link} from 'svelte-routing';

    export let message;
</script>

<section class="BoxSc">
    <div class="CardBoxWrap">
        <div class="CardBoxContainer">
            <div class="CardBody">
                <div class="CardTitle">
                    <div><h5>{$_('msg_edit_validator')}</h5></div>
                </div>
            </div>
            <div class="CardList">
                <table class="CardTable MsSend">
                    <thead>
                    </thead>
                    <tbody>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('moniker')}</p>
                        </td>

                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark MSG">{message.description.moniker}</p>
                        </td>
                    </tr>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('website')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark MSG">{message.description.website}</p>
                        </td>
                    </tr>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('security_contact')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark MSG">{message.description.security_contact}</p>
                        </td>
                    </tr>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('details')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark MSG">{message.description.details}</p>
                        </td>
                    </tr>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('validator_address')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark MSG">
                                <Link to="/account/{message.validator_address}">
                                    {message.validator_address}
                                </Link>
                            </p>
                        </td>
                    </tr>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('commission_rate')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark">
                                {message.commission_rate}
                            </p>
                        </td>
                    </tr>
                    <tr>
                        <td class="px-2 py-1 TextAL">
                            <p class="FS15 fwM Color_Dark">{$_('min_self_delegation')}</p>
                        </td>
                        <td class="px-2 py-1 TextAR">
                            <p class="FS15 fwM Color_Dark">
                                {message.min_self_delegation}
                            </p>
                        </td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</section>