# ✨🎨 AstroCanvas ✨🎨

Frontend repository: [astro-canvas-frontend](https://github.com/chainapsis/astro-canvas-frontend)  
Custom Keplr extension for hackathon (preliminary IBC support): [Keplr-extension](https://github.com/chainapsis/keplr-extension/tree/hackaton) [Release](https://github.com/chainapsis/keplr-extension/releases/tag/v0.6.0-hackathon)  
Website: [hackathon.keplr.app](https://hackathon.keplr.app/)  

### Decentralize Staking, Colorfully 🌈

### Introduction

Stake decetralization, or the spread of voting power across multiple validators, is one crucial element to increase the robustness and decentralization of a public proo-of-stake blockchain. As voting power gets concentrated to the hands of a few powerful validators, number of colluding entities needed to create byzantine behavior reduces.

![top-6-validators](img/hub-validators.jpeg)

As for the Cosmos Hub, only the top 6 validators need to collude in order to amass 33.4% of voting power–which can effectively halt the blockchain. While whether such actions are realistic may be questionable, it's no doubt that there is room for improvement.

While ideas such as variable staking reward rate, variable slashing rate, and more have been proposed, such ideas take significant economic research and security analysis to ensure the cryptoeconomics of the blockchain remains robust.

### Paint the picture of how you want **your** Cosmos to looks like

> So what if we could incentivize stake decentralization without changing the fundamentals of the blockchain cryptoeconomics, by using  **personal amusement**?

Astrocanvas is a game and a radical social experimentation in stake decentralization. We use elements of entertainment, scarcity, competition, and economics to incentivize voting power distribution of a proof-of-stake blockchain.

![place-reddit](img/place-reddit.png)

We were inspired by the 2017 Reddit april fools project 'Place'. Where each account was allowed to change the color of one pixel on a 1000x1000 pixel canvas–**but only every 5 minutes**. In just 72 hours, over 1 million users participated. During that time, factions were formed, betrayals happened, virtual wars were fought, and ultimately a picture was drawn that represented the community of people that drew it. Place was an experimentation of digital scarcity, human psychology, and competition.

In Astrocanvas, delegators are given specific `colorToken` that represents the right to change one pixel in the canvas when they delegate their staking token to a Hub validator. 

The catch? **Not all `colorToken` are the same.**

Delegators can earn `colorToken` of a specific color (white, black, red, etc) depending on the voting power of the delegated validator. So for example, validator with *#1 to #10* rank in voting power gives delegators `colorTokenWhite` which only allows you to place a white pixel on the canvas. If you want to place a blue pixel, you need `colorTokenBlue` which you may only receive when you delegate to a validator with voting power ranging from *#80 to #90*.

The idea here is as following:
* People want to draw specific 'art' on the canvas, which would require a wide set of colors to do. Which will disincentivize staking to a high voting power validator in order to acquire other colors.
* The picture drawn on the canvas, and the dominance of a specific color will accurately reflect the community of delegators and the state of the hub.
* Validators will incentivize drawing of specific images that are dominant in the color of their respective `colorToken` in order to bring in delegations.
* Potentially many factions and sub-communities will form in order to own a piece of the digitally scarce, immutable, decentralized canvas.

### A few quirks and features

* The `colorToken` is burnt when it's used to change the pixel on the Astrocanvas
* The amount of `colorToken`s you receive is relative to the amount you stake
* You receive new `colorToken`s every `n` blocks
* You can only use `colorToken`s every `n` blocks
* More in progress...

### How we built it

We have used an early implementation of [ICS27(Interchain Accounts)](https://github.com/cosmos/ics/tree/master/spec/ics-027-interchain-accounts) for interchain staking to support staking to another chain. In the current implementation, the interchain staking module allows the foreign chain to mint the representative color tokens that include information such as the source port, source channel, validator address once the staking is successful.

The AstroCanvas zone is built using a modified Cosmos-SDK, and records the state of the canvas (and each of the pixels) on the ledgers. The canvas module manages the canvas and the specific pixels of the canvas. Any user can create the canvas with the canvas id, x-coordinate, y-coordinate, allowed denoms, etc. Users can point to a specific canvas id by using the assets that the canvas permissions. The front-end draws the denoms that exist on the canvas, and assigns specific colors for the on-chain information. Our current implementation works on the optimistic assumption that 'all IBC packets will succeed' as the current relayer implementation has limitations in sending acknoledgements packets--and we didn't have time to revise the code for that part.

## How to use AstroCanvas

**How to install Keplr IBC**

1. Download [Keplr IBC (Hackathon Release)](https://github.com/chainapsis/keplr-extension/releases/tag/v0.6.0-hackathon) and unzip/extract the files

2. Open Chrome and click on the menu bar (the three dots on the top right corner of the browser). Go to 'More tools' -> 'Extensions'

3. Turn 'Developer Mode' on (top right corner of the page)

4. Click on 'Load Unpacked', and choose the Keplr IBC folder you downloaded in step 1.

5. Open Keplr Extension and create/import an account.

**How to use AstroCanvas**

1. Get some tokens from the Astro Hub and Astro Zone faucet (50 tokens per 5 minutes)

2. Send the Astro Hub tokens to Astro Zone via IBC send (Note: Astro Hub token is the staking token that's used to generate the `colorTokens` from within the Astro Zone. Astro Zone's tokens are only used as gas tokens)

3. Go to [hackathon.keplr.app](https://hackathon.keplr.app) and click on 'Register'

4. Delegate to validators and receive specific `colorTokens` __(Note: Your color token balance regenerate every 5 mins. i.e. if you have spent all 50 `whiteTokens` that you initially received, your `whiteTokens` be 50 again after 5 minutes since you've last used the `whiteToken` on the canvas)__

5. Place pixels on the canvas, and create your drawing! (We suggest you don't place more than 30 pixels at a time as too many pixels can cause issues with transactions failing)
