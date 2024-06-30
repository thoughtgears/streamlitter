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

const GIT_SHA = childProcess.exec('git rev-parse --short HEAD', (err, stdout, stderr) => {
    if (err) {
        console.error(`Error getting Git SHA: ${stderr}`);
        process.exit(1);
    }
    console.log(`Current Git SHA: ${stdout.trim()}`);
});


function chooseBinary() {
    const platform = os.platform()
    const arch = os.arch()

    if (platform === 'linux' && arch === 'x64') {
        return `main-linux-amd64-${VERSION}-${GIT_SHA}`
    }
    if (platform === 'linux' && arch === 'arm64') {
        return `main-linux-arm64-${VERSION}-${GIT_SHA}`
    }
    if (platform === 'darwin' && arch === 'x64') {
        return `main-darwin-amd64-${VERSION}-${GIT_SHA}`
    }
    if (platform === 'darwin' && arch === 'arm64') {
        return `main-darwin-arm64-${VERSION}-${GIT_SHA}`
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
