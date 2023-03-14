export class Token {
  constructor({denom, symbol, amount, precision}) {
    this.denom = network.networkToken.base;
    this.symbol = network.networkToken.symbol;
    this.precision = network.networkToken.precision;

    if (denom !== undefined) {
      this.denom = denom;
    }
    if (symbol !== undefined) {
      this.symbol = symbol;
    }
    if ( precision !== undefined) {
      this.precision = precision;
    }
    this.amount = amount;
  }

  string() {
    let amount = BigInt(this.amount);
    if (this.symbol === undefined) {
      return `${amount}${this.symbol}`
    }

    var factor = BigInt(1)

    for (let i = 0; i < this.precision; i++) {
      factor *= BigInt(10);
    }

    let str = this.amount.toString().padStart(this.precision+1, "0") // +1 for the leading zero

    let integer = str.substring(0, str.length - this.precision);
    // TODO remove trailing zero
    let frac = str.substring(str.length - this.precision).replace(/0+$/,'') || "00";
    

    return `${integer}.${frac} ${this.symbol.toUpperCase()}`;

  }
}
