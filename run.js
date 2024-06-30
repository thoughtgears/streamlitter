const childProcess = require('child_process');
const os = require('os')
const process = require('process')
const fs = require('fs')
const path = require('path')


const versionFilePath = path.join(__dirname, '.version');
let VERSION;
try {
    VERSION = fs.readFileSync(versionFilePath, 'utf8').trim();
} catch (err) {
    console.error(`Error reading version file: ${err.message}`);
    process.exit(1);
}

function chooseBinary() {
    const platform = os.platform()
    const arch = os.arch()

    if (platform === 'linux' && arch === 'x64') {
        return `main-linux-amd64-${VERSION}`
    }
    if (platform === 'linux' && arch === 'arm64') {
        return `main-linux-arm64-${VERSION}`
    }
    if (platform === 'darwin' && arch === 'x64') {
        return `main-darwin-amd64-${VERSION}`
    }
    if (platform === 'darwin' && arch === 'arm64') {
        return `main-darwin-arm64-${VERSION}`
    }

    console.error(`Unsupported platform (${platform}) and architecture (${arch})`)
    process.exit(1)
}

function main() {
    const binary = chooseBinary()
    console.log(`Using binary: ${binary}`)
    const mainScript = `${__dirname}/bin/${binary}`
    const spawnSyncReturns = childProcess.spawnSync(mainScript, { stdio: 'inherit' })
    const status = spawnSyncReturns.status
    if (typeof status === 'number') {
        process.exit(status)
    }
    process.exit(1)
}

if (require.main === module) {
    main()
}
