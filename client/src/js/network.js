export function setNetwork(net) {
  try {
    if (net === undefined || net === '') {
      net = 'testnet';
    } else {
      if (!config.network.hasOwnProperty(net)) {
        net = 'testnet';
      }
    }
  } catch (e) {
    net = 'testnet';
  } finally {
    console.log(`current network: ${net}`);
    console.log(config.network[net]);
    window.network = config.network[net];
    localStorage.setItem('network', net);
  }
}

export function getNetwork() {
  let network = localStorage.getItem('network');
  if (network === 'undefined' || network === undefined) {
    network = 'testnet';
  }

  // console.debug(`current network : ${network}`);
  return network;
}
