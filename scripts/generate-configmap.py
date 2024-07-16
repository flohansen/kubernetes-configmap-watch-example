import glob
import json

if __name__ == '__main__':
    with open('deployment/kubernetes/importer.configmap.yaml', 'w') as f:
        f.write('\n'.join([
                "apiVersion: v1",
                "kind: ConfigMap",
                "metadata:",
                "  name: test-product-data",
                "data:",
                "  products.json: |",
                "",
                ]))

        contents = []

        for file in glob.glob('**/*.json'):
            with open(file) as f2:
                contents.append(json.load(f2))

        dump = [' ' * 4 + line for line in json.dumps(contents, indent=2).splitlines()]
        f.write('\n'.join(dump))
