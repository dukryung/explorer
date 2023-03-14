<script>
    import { Link } from 'svelte-routing';
    import { _, init, locales } from 'svelte-i18n';
    import AppLogo from '../../assets/AppLogoSvg.svelte';
    import NiktoNetText from '../../assets/NiktoNetText.svelte';
    import IconCurrentDot from '../../assets/IconCurrentDot.svelte';
    import IconNomalDot from '../../assets/IconNomalDot.svelte';
    import IconDownArrow from '../../assets/IconDownArrow.svelte';
    import { getNetwork, setNetwork } from 'js/network';
    import { onDestroy, onMount } from 'svelte';
    import { networkConnection, localeGlobalState } from 'js';

    let dropdownMenuState = false;
    let networkState = false;
    let langMenuState = false;
    let currentLocale;
    let currentLocales = [];
    const unsubscribe = networkConnection.subscribe(
        conn => (networkState = conn)
    );

    const menus = [
        { to: '/blocks', text: 'blocks' },
        { to: '/txs', text: 'transactions' },
        { to: '/tokens', text: 'tokens' },
        { to: '/validators', text: 'validators' },
        { to: '/nfts', text: 'nfts' }
    ];
    const networks = [];
    getNetworks();

    onMount(() => {
        window.addEventListener('resize', () => {
            if (
                document.getElementById('Header').classList.contains('nav-open')
            ) {
                if (window.screen.width > 960) {
                    closeMobileMenu();
                }
            }
        });
    });

    onDestroy(() => {});

    function onMobileMenu() {
        document.getElementById('Header').classList.toggle('nav-open');
        dropdownMenuState = !dropdownMenuState;
    }

    function closeMobileMenu() {
        document.getElementById('Header').classList.remove('nav-open');
        dropdownMenuState = false;
    }

    const onDropdownMenu = () => {
        dropdownMenuState = true;
    };

    const closeDropdownMenu = () => {
        if (window.screen.width < 960) {
            return;
        }
        dropdownMenuState = false;
    };

    function getNetworks() {
        for (const network in config.network) {
            networks.push({
                key: network,
                data: config.network[network]
            });
        }
    }

    function switchNetwork(network) {
        setNetwork(network);
        window.location = '/';
    }
    const toggleLangMenu = () => {
        langMenuState = !langMenuState;
    };
    localeGlobalState.subscribe(value => (currentLocale = value));
    $: if (currentLocale !== undefined) {
        init({
            fallbackLocale: 'en-US',
            initialLocale: currentLocale
        });

        currentLocales = [];
        for (const locale of $locales) {
            if (locale === currentLocale) {
                currentLocales.splice(0, 0, locale);
            } else {
                currentLocales.push(locale);
            }
        }

        currentLocales = currentLocales;
    }
    const onChangeLocale = locale => {
        localeGlobalState.set(locale);
        localeGlobalState.update(v => (locale = v));
            toggleLangMenu();
    };
</script>

{#if !networkState}
    <div class="network">
        {$_('network_error')}
    </div>
{/if}
<header class="cd-morph-dropdown" id="Header">
    <!-- TODO insert EXAIS logo
    <Link class="logo" to="/">
        <AppLogo fill="white"/>
    </Link>
    -->
    <a class="nav-trigger" on:click={onMobileMenu}
        >Open Nav<span aria-hidden="true" /></a
    >
    <nav class="main-nav">
        <ul>
            {#each menus as menu}
                <li class="">
                    <Link to={menu.to}>
                        {$_(menu.text)}
                    </Link>
                </li>
            {/each}
        </ul>
    </nav>
    <div
        class="HeaderLangButton"
        on:mouseenter={toggleLangMenu}
        on:mouseleave={toggleLangMenu}
    >
        {$_(currentLocale)}
        {#if langMenuState}
            <div class="HeaderLangWrap">
                {#each currentLocales as locale, i}
                    {#if currentLocale === locale}
                        <div
                            class="LangButton selected"
                            on:click={() => {
                                onChangeLocale(locale);
                            }}
                        >
                            {$_(locale)}
                        </div>
                    {:else}
                        <div
                            class="LangButton unselected"
                            on:click={() => {
                                onChangeLocale(locale);
                            }}
                        >
                            {$_(locale)}
                        </div>
                    {/if}
                {/each}
            </div>
        {/if}
    </div>

    <!--            <li class="has-dropdown button" data-content="Net"-->
    <!--                on:mouseenter={onDropdownMenu}-->
    <!--                on:mouseleave={closeDropdownMenu}-->
    <!--            >-->
    <!--                <a class="BtnDIm">-->
    <!--                    Service-->
    <!--                    <span class="Darrow">-->
    <!--                        <IconDownArrow/>-->
    <!--                    </span>-->
    <!--                </a>-->
    <!--            </li>-->

    <div class="morph-dropdown-wrapper">
        <div class="dropdown-list">
            <ul>
                {#each menus as menu}
                    <li class="">
                        <Link
                            class="dep1"
                            to={menu.to}
                            on:click={closeMobileMenu}
                        >
                            {$_(menu.text)}
                        </Link>
                    </li>
                {/each}
                <!--                <li id="Net" class="dropdown Net"-->
                <!--                    on:mouseenter={onDropdownMenu}-->
                <!--                    on:mouseleave={closeDropdownMenu}-->
                <!--                >-->
                <!--                    <a class="dep1">Service</a>-->
                <!--                    {#if dropdownMenuState}-->
                <!--                        <div class="content">-->
                <!--                            <ul>-->
                <!--                                {#each networks as network}-->
                <!--                                    {#if network.key === getNetwork()}-->
                <!--                                        <li>-->
                <!--                                            <a class="CurrentBtnDim">-->
                <!--                                                <p class="CurrentTextColor">-->
                <!--                                                    <NiktoNetText/>-->
                <!--                                                    <spa>{network.data.name}</spa>-->
                <!--                                                    <IconCurrentDot/>-->
                <!--                                                </p>-->
                <!--                                            </a>-->
                <!--                                        </li>-->
                <!--                                    {:else}-->
                <!--                                        <li>-->
                <!--                                            <a class="CurrentBtn" on:click={()=>{-->
                <!--                                                    switchNetwork(network.key);-->
                <!--                                                }}>-->
                <!--                                                <p class="TextNomal">-->
                <!--                                                    <NiktoNetText/>-->
                <!--                                                    <spa>{network.data.name}</spa>-->
                <!--                                                    <IconNomalDot/>-->
                <!--                                                </p>-->
                <!--                                            </a>-->
                <!--                                        </li>-->
                <!--                                    {/if}-->
                <!--                                {/each}-->
                <!--                            </ul>-->
                <!--                        </div>-->
                <!--                    {/if}-->
                <!--                </li>-->
            </ul>
            <div class="HeaderMobileLangButton">
                Language
                <div class="LangButton">
                    {#each $locales as locale, i}
                        {#if currentLocale === locale}
                            <a
                                class="selected"
                                on:click={closeMobileMenu}
                                on:click={() => {
                                    onChangeLocale(locale);
                                }}>{$_(locale)}</a
                            >
                        {:else}
                            <a
                                class="unselected"
                                on:click={closeMobileMenu}
                                on:click={() => {
                                    onChangeLocale(locale);
                                }}>{$_(locale)}</a
                            >
                        {/if}

                        {#if i + 1 !== $locales.length}
                            <span>&#183;</span>
                        {/if}
                    {/each}
                </div>
            </div>
            <div class="bg-layer" aria-hidden="true" />
        </div>
    </div>
</header>

<style>
    .network {
        position: fixed;
        z-index: 999;
        width: 100%;
        height: 30px;
        background-color: red;
        color: white;
        text-align: center;
    }
</style>
