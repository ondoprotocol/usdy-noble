cd proto
buf generate --template buf.gen.gogo.yaml
cd ..

cp -r github.com/ondoprotocol/aura/* ./
rm -rf github.com
