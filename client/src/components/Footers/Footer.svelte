<script>
    import { _ } from 'svelte-i18n';
    import { Link } from 'svelte-routing';
    import HesseggLogoSvg from '../../assets/HesseggLogoSvg.svelte';
    import HesseggLogoTextSvg from '../../assets/HesseggLogoTextSvg.svelte';
    import IconEmail from '../../assets/IconEmail.svelte';
</script>

<footer class="footer">
    <div class="FooterWrap">
        <div class="Fmenu">
            <a class="BtnEmail" href="mailto://info@exais.net">
                <IconEmail />
            </a>
        </div>
        <!--
        <div class="CPleft">
            <HesseggLogoSvg/>
            <HesseggLogoTextSvg/>
            <p class="FCR FS12">Copyright © HESSEGG. All Rights Reserved.</p>
        </div>
-->
        <!--<div class="TexTMenu"> <Link class="ToU" to="">Terms of Use</Link><Link class="PP" to="">Privacy Policy</Link></div>-->
    </div>
</footer>
