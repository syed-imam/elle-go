(ns elle-bridge
  (:require [elle.list-append :as la]
            [jepsen.history :as h]
            [clojure.edn :as edn])
  (:import [java.io PushbackReader FileReader]))

(defn -main [file]
  (with-open [r (PushbackReader. (FileReader. file))]
    (let [ops    (doall (take-while some? (repeatedly #(edn/read {:eof nil} r))))
          hist   (h/history ops)
          result (la/check {:consistency-models [:serializable]} hist)]
      (prn {:valid?        (:valid? result)
            :anomaly-types (vec (:anomaly-types result))}))))
