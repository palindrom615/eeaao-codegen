export const swaggerPetstoreApi = {
        listPets: (data) => {
            return fetch("/pets", {
                method: "get",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            })
        },
        createPets: (data) => {
            return fetch("/pets", {
                method: "post",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            })
        },
        showPetById: (data) => {
            return fetch("/pets/{petId}", {
                method: "get",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            })
        },
}
