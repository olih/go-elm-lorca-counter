port module Ports exposing (sendIncMessage, activeCounter)

port sendIncMessage : Int -> Cmd msg

port activeCounter : (Int -> msg) -> Sub msg