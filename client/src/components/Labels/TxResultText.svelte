<script>
    import {_} from 'svelte-i18n';
    import IconSmSuccess from '../../assets/IconSmSuccess.svelte';
    import IconSmFailed from '../../assets/IconSmFailed.svelte';

    export let code = 0;
</script>

{#if code === 0}
    <!-- success -->
    {$_('success')}<IconSmSuccess/>
{:else }
    <!-- failed -->
    {$_('failed')}<IconSmFailed/>
{/if}