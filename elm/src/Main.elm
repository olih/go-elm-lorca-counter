import Browser
import Html exposing (Html, button, div, text)
import Html.Events exposing (onClick)
import Html.Attributes as Attr
import Ports exposing (sendIncMessage, activeCounter)


-- MODEL

type alias Model =
    { counter : Int
  }

initialModel : Model
initialModel = {
  counter = 0
 }

-- UPDATE

type AppMsg = Increment | Decrement | RefreshCounter Int

update : AppMsg -> Model -> ( Model, Cmd AppMsg )
update msg model =
  case msg of
    Increment ->
      (model, sendIncMessage 1)

    Decrement ->
      (model, sendIncMessage -1)

    RefreshCounter counter ->
      ( { model | counter = counter }, Cmd.none )

-- VIEW

view : Model -> Html AppMsg
view model =
  div []
    [ button [ Attr.class "button", onClick Decrement ] [ text "-" ]
    , div [] [ text (String.fromInt model.counter) ]
    , button [ Attr.class "button is-danger", onClick Increment ] [ text "+" ]
    ]

subscriptions : Model -> Sub AppMsg
subscriptions model =
    activeCounter RefreshCounter


type alias Flags =
    {}


init : Flags -> ( Model, Cmd AppMsg )
init flags =
    ( initialModel, Cmd.none )


main : Program Flags Model AppMsg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }
