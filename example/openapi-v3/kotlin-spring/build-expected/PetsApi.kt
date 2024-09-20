package com.example.api

import io.swagger.v3.oas.annotations.ExternalDocumentation
import io.swagger.v3.oas.annotations.Operation
import io.swagger.v3.oas.annotations.Parameter
import io.swagger.v3.oas.annotations.enums.ParameterIn
import io.swagger.v3.oas.annotations.headers.Header
import io.swagger.v3.oas.annotations.media.ArraySchema
import io.swagger.v3.oas.annotations.media.Content
import io.swagger.v3.oas.annotations.media.Schema
import io.swagger.v3.oas.annotations.parameters.RequestBody
import io.swagger.v3.oas.annotations.responses.ApiResponse
import io.swagger.v3.oas.annotations.tags.Tag
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestParam
import javax.annotation.Generated


@Generated
@Tag(
    name = "pets",
    description = "Everything about your Pets",
    externalDocs = ExternalDocumentation(description = "Find out more", url = "http://swagger.io")
)
interface PetsApi {
    @Operation(
        summary = "List all pets",
        operationId = "listPets",
        responses = [
            ApiResponse(
                responseCode = "200",
                description = "A list of pets",
                content = [
                    Content(
                        mediaType = "application/json",
                        array = ArraySchema(
                            arraySchema = Schema(name = "Pets"),
                            schema = Schema(implementation = Pet::class),
                            maxItems = 100
                        )
                    )
                ],
                headers = [
                    Header(
                        name = "x-next",
                        description = "A link to the next page of responses",
                        schema = Schema(type = "string")
                    )]
            ),
            ApiResponse(
                responseCode = "default",
                description = "unexpected error",
                content = [Content(mediaType = "application/json", schema = Schema(implementation = Error::class))]
            )
        ]
    )
    @GetMapping("/pets")
    fun listPets(
        @Parameter(
            `in` = ParameterIn.QUERY,
            description = "How many items to return at one time (max 100)",
            required = false,
            schema = Schema(type = "integer", maximum = "100", format = "int32")
        )
        @RequestParam("limit")
        limit: Int?,
    ): List<Pet>

    @Operation(
        summary = "Create a pet",
        operationId = "createPets",
        responses = [
            ApiResponse(
                responseCode = "201",
                description = "Null response"
            ),
            ApiResponse(
                responseCode = "default",
                description = "unexpected error",
                content = [Content(mediaType = "application/json", schema = Schema(implementation = Error::class))]
            )
        ]
    )
    @PostMapping("/pets")
    fun createPets(
        @RequestBody(
            description = "Pet to add to the store",
            required = true,
            content = [
                Content(
                    mediaType = "application/json",
                    schema = Schema(implementation = Pet::class)
                )
            ]
        )
        @org.springframework.web.bind.annotation.RequestBody
        pet: Pet,
    )

    @Operation(
        summary = "Info for a specific pet",
        operationId = "showPetById",
        responses = [
            ApiResponse(
                responseCode = "200",
                description = "Expected response to a valid request",
                content = [Content(mediaType = "application/json", schema = Schema(implementation = Pet::class))]
            ),
            ApiResponse(
                responseCode = "default",
                description = "unexpected error",
                content = [Content(mediaType = "application/json", schema = Schema(implementation = Error::class))]
            )
        ]
    )
    @GetMapping("/pets/{petId}")
    fun showPetById(
        @Parameter(
            `in` = ParameterIn.PATH,
            description = "The id of the pet to retrieve",
            required = true,
            schema = Schema(type = "string")
        )
        petId: Long,
    ): Pet
}

@Generated
@Schema(
    name = "Pet",
    requiredProperties = ["id", "name"],
)
data class Pet(
    val id: Long,
    val name: String,
    val tag: String?,
)

@Generated
@Schema(
    name = "Error",
    requiredProperties = ["code", "message"],
)
data class Error(
    val code: Int,
    val message: String,
)
