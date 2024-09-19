def main():
    specs = eeaao_codegen.loadSpecsGlob("proto", "*.proto")
    print(specs)

    specs = eeaao_codegen.loadSpecsGlob("json", "*.json")
    v2Spec = specs["petstorev2.json"]
    info = v2Spec["info"]
    tags = v2Spec["tags"]
    paths = v2Spec["paths"]
    # print(spec["components"])
    print(v2Spec)

    v3Spec = specs["petstore.json"]
    print(v3Spec)

    trainTravelSpec = specs["traintravel.json"]
    print(trainTravelSpec)
    for k in v2Spec:
        print(k)
        print(v2Spec[k])
        print("\n")
