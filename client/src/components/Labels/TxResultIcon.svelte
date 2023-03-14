<script>
    import {_} from 'svelte-i18n';
    import IconBigSuccess from '../../assets/IconBigSuccess.svelte';
    import IconBigFailed from '../../assets/IconBigFailed.svelte';

    export let code = 0;
</script>

{#if code === 0}
    <!-- success -->
    <div class="ImgWrap">
        <IconBigSuccess/>
    </div>
{:else }
    <!-- failed -->
    <div class="ImgWrap">
        <IconBigFailed/>
    </div>
{/if}