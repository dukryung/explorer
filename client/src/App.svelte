<script>
    import { Route, Router } from 'svelte-routing';
    import Header from 'Components/Headers/Header.svelte';
    import Blocks from 'Layouts/Blocks.svelte';
    import Dashboard from 'Layouts/Dashboard.svelte';
    import Transactions from 'Layouts/Transactions.svelte';
    import Block from './layouts/Block.svelte';
    import Transaction from './layouts/Transaction.svelte';
    import Footer from './components/Footers/Footer.svelte';
    import Account from './layouts/Account.svelte';
    import Tokens from './layouts/Tokens.svelte';
    import Validator from 'Layouts/Validator';
    import NotFound from 'Layouts/NotFound';
    import Validators from 'Layouts/Validators';
    import NFToken from './layouts/NFToken.svelte';
    import NFTokenFixed from './layouts/NFTokenFixed.svelte';
    import NFTokens from './layouts/NFTokens.svelte';
    import NFTCollection from './layouts/NFTCollection.svelte';
    const url = '';
</script>

<div class="wrap">
    <Router {url}>
        <Header />
        <Route path="/">
            <Dashboard />
        </Route>
        <Route path="/account/:address" let:params>
            <Account address={params.address} />
        </Route>
        <Route path="/validator/:address" let:params>
            <Validator address={params.address} />
        </Route>
        <Route path="/validators">
            <Validators />
        </Route>
        <Route path="/block/:height" let:params>
            <Block height={params.height} />
        </Route>
        <Route path="/blocks">
            <Blocks />
        </Route>
        <Route path="/tx/:hash" let:params>
            <Transaction hash={params.hash} />
        </Route>
        <Route path="/txs">
            <Transactions />
        </Route>
        <Route path="/tokens">
            <Tokens />
        </Route>
        <Route path="/nfts">
            <NFTokens />
        </Route>
        <Route path="/nfts/collection/:col_id" let:params>
            <NFTCollection col_id={params.col_id} />
        </Route>
        <Route path="/nft/:token_id/" let:params>
            <NFTokenFixed token_id={params.token_id} />
        </Route>
        <Route path="/nft">
            <NFToken />
        </Route>
        <Route path="*">
            <NotFound />
        </Route>
    </Router>
    <div class="bgImage" />
</div>
<Footer />

<style global lang="postcss">
</style>
