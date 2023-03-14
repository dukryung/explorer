<script>
    import MsgSend from './Types/MsgSend.svelte';
    import MsgIssue from './Types/MsgIssue.svelte';
    import MsgMint from './Types/MsgMint.svelte';
    import MsgBurn from './Types/MsgBurn.svelte';
    import MsgName from './Types/MsgName.svelte';
    import MsgTransfer from './Types/MsgTransfer.svelte';
    import Default from 'Components/Messages/Types/Default.svelte';
    import MsgUnjail from './Types/MsgUnjail.svelte';
    import MsgDelegate from './Types/MsgDelegate.svelte';
    import MsgBeginRedelegate from './Types/MsgBeginRedelegate.svelte';
    import MsgUndelegate from './Types/MsgUndelegate.svelte';
    import MsgEditValidator from './Types/MsgEditValidator.svelte';
    import MsgConvert from './Types/MsgConvert.svelte';
    import MsgTransferNFT from './Types/MsgTransferNFT.svelte';
    export let messages = [];

    function messageComponent(message) {
      switch (message['@type']) {
        case '/cosmos.bank.v1beta1.MsgSend':
          return MsgSend;
          case '/nikto.bankz.v1.MsgTransfer':

              return MsgTransfer;
          case '/nikto.name.v1.MsgRegisterName':
               return MsgName;
        case '/klaatoo.bankz.v1beta1.MsgIssue':
          return MsgIssue;
        case '/klaatoo.bankz.v1beta1.MsgMint':
          return MsgMint;
        case '/klaatoo.bankz.v1beta1.MsgBurn':
          return MsgBurn;
        case '/cosmos.slashing.v1beta1.MsgUnjail':
          return MsgUnjail;
        case '/cosmos.staking.v1beta1.MsgBeginRedelegate':
          return MsgBeginRedelegate;
        case '/cosmos.staking.v1beta1.MsgUndelegate':
          return MsgUndelegate;
        case '/cosmos.staking.v1beta1.MsgDelegate':
          return MsgDelegate;
          case '/cosmos.staking.v1beta1.MsgEditValidator':
          return MsgEditValidator;
        case '/nikto.convert.v1.MsgConvert':
          return MsgConvert;
          case '/nikto.nft.v1.MsgTransferNFT':
            return MsgTransferNFT;
        default:
          return Default;
      }
    }
</script>

{#each messages as message}
    <svelte:component this={messageComponent(message)} message={message}/>
{/each}



