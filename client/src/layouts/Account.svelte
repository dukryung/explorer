<script>
    import { _ } from 'svelte-i18n';
    import Meta from 'svelte-meta';
    import AccountSummaryCard from '../components/Cards/AccountSummaryCard.svelte';
    import TxListCard from '../components/Cards/TxListCard.svelte';
    import SearchCard from '../components/Cards/SearchCard.svelte';
    import AccountTokenListCard from 'Components/Cards/AccountTokenListCard';
    import DelegationCard from '../components/Cards/DelegationCard.svelte';
    import RedelegationCard from '../components/Cards/RedelegationCard.svelte';
    import UnbondingCard from '../components/Cards/UnbondingCard.svelte';
    import NftCollectionCard from '../components/Cards/NFTCollectionCard.svelte';

    export let address;

    const txLimit = 10;
    const txPage = 0;

    const tokenPage = 0;
    const tokenLimit = 12;
</script>

<Meta title="{$_('meta_title')} | {$_('account')} - {address}" />

<SearchCard title={$_('account_details')} />

<AccountSummaryCard {address} />

<AccountTokenListCard {address} />
<!--
<DelegationCard
        address={address}
/>

<RedelegationCard
        address={address}
/>

<UnbondingCard
        address={address}
/>
-->
<TxListCard
    type={klaatoo.SUBSCRIBE}
    method={'tx.address'}
    {address}
    {txLimit}
    {txPage}
/>
<NftCollectionCard
    method={'nft.listbyowner'}
    {address}
    {tokenPage}
    {tokenLimit}
/>
<NftCollectionCard
    method={'nft.listbyissuer'}
    {address}
    {tokenPage}
    {tokenLimit}
/>
