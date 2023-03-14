<script>
    import SearchCard from '../components/Cards/SearchCard.svelte';
    import {_} from 'svelte-i18n';
    import Meta from 'svelte-meta';
    import TokenListCard from '../components/Cards/TokenListCard.svelte';
</script>

<Meta
        title="{$_('meta_title')} | {$_('tokens')}"
        description="{$_('meta_description')}"
/>

<SearchCard title={$_('tokens')}/>
<TokenListCard/>
