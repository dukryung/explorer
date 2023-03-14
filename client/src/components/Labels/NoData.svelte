<script>
    import {_} from 'svelte-i18n';

    export let description = $_('no_data');
</script>

<div class="NoData">
    <div class="NoDataWrap">
        <img src="/images/nodata.svg">
        <p>{description}</p>
    </div>
</div>