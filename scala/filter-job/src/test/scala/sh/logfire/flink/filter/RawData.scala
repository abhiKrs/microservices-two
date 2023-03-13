package sh.logfire.flink.filter

import io.circe._
import io.circe.generic.semiauto.deriveDecoder
import io.circe.parser._
import org.apache.avro.Schema
import org.apache.avro.generic.GenericData

object RawData {
  case class RawRecord(
                        platform: String,
                        message: String,
                        msgid: String,
                        pid: Int,
                        level: String,
                        dt: String,
                        version: Int,
                    )

  case class Data(data: Seq[RawRecord])
  implicit val rawRecordDecoder: Decoder[RawRecord] = deriveDecoder[RawRecord]

  implicit val dataDecoder: Decoder[Data] = deriveDecoder[Data]



  def getGenericRecords(): Seq[GenericData.Record] = {
    val jsonRecordsString = """
      |{
      |    "data" : [
      |        {"platform":"docker","message":"You're not gonna believe what just happened","msgid":"ID841","pid":9447,"level":"debug","dt":"2023-02-15T07:15:03.099Z","version":1},
      |        {"platform":"kubernetes","message":"Great Scott! We're never gonna reach 88 mph with the flux capacitor in its current state!","msgid":"ID840","pid":8196,"level":"debug","dt":"2023-02-15T07:15:04.098Z","version":1},
      |        {"platform":"nginx","message":"Pretty pretty pretty good","msgid":"ID584","pid":9715,"level":"debug","dt":"2023-02-15T07:15:05.098Z","version":1},
      |        {"platform":"ubuntu","message":"#hugops to everyone who has to deal with this","msgid":"ID77","pid":1273,"level":"err","dt":"2023-02-15T07:15:06.098Z","version":2},
      |        {"platform":"docker","message":"There's a breach in the warp core, captain","msgid":"ID383","pid":5524,"level":"info","dt":"2023-02-15T07:15:07.098Z","version":1},
      |        {"platform":"heroku","message":"Pretty pretty pretty good","msgid":"ID517","pid":3682,"level":"notice","dt":"2023-02-15T07:15:08.098Z","version":1},
      |        {"platform":"kubernetes","message":"There's a breach in the warp core, captain","msgid":"ID142","pid":7632,"level":"info","dt":"2023-02-15T07:15:09.098Z","version":1},
      |        {"platform":"kubernetes","message":"There's a breach in the warp core, captain","msgid":"ID932","pid":494,"level":"alert","dt":"2023-02-15T07:15:10.098Z","version":2},
      |        {"platform":"aws","message":"#hugops to everyone who has to deal with this","msgid":"ID189","pid":3210,"level":"warning","dt":"2023-02-15T07:15:11.098Z","version":1},
      |        {"platform":"gcp","message":"Great Scott! We're never gonna reach 88 mph with the flux capacitor in its current state!","msgid":"ID720","pid":7176,"level":"notice","dt":"2023-02-15T07:15:12.098Z","version":2},
      |        {"platform":"kubernetes","message":"We're gonna need a bigger boat","msgid":"ID819","pid":9134,"level":"notice","dt":"2023-02-15T07:15:13.098Z","version":1},
      |        {"platform":"nginx","message":"A bug was encountered but not in Vector, which doesn't have bugs","msgid":"ID214","pid":1950,"level":"crit","dt":"2023-02-15T07:15:14.098Z","version":1},
      |        {"platform":"nginx","message":"There's a breach in the warp core, captain","msgid":"ID862","pid":7354,"level":"warning","dt":"2023-02-15T07:15:15.098Z","version":1},
      |        {"platform":"heroku","message":"Pretty pretty pretty good","msgid":"ID878","pid":1738,"level":"notice","dt":"2023-02-15T07:15:16.098Z","version":1},
      |        {"platform":"docker","message":"Take a breath, let it go, walk away","msgid":"ID516","pid":8131,"level":"err","dt":"2023-02-15T07:15:17.098Z","version":1},
      |        {"platform":"nginx","message":"You're not gonna believe what just happened","msgid":"ID211","pid":6733,"level":"warning","dt":"2023-02-15T07:15:18.099Z","version":1},
      |        {"platform":"nginx","message":"Maybe we just shouldn't use computers","msgid":"ID320","pid":6046,"level":"crit","dt":"2023-02-15T07:15:19.098Z","version":2},
      |        {"platform":"gcp","message":"Great Scott! We're never gonna reach 88 mph with the flux capacitor in its current state!","msgid":"ID947","pid":220,"level":"info","dt":"2023-02-15T07:15:20.099Z","version":2},
      |        {"platform":"gcp","message":"Pretty pretty pretty good","msgid":"ID239","pid":6472,"level":"crit","dt":"2023-02-15T07:15:21.098Z","version":2},
      |        {"platform":"heroku","message":"Pretty pretty pretty good","msgid":"ID490","pid":4576,"level":"info","dt":"2023-02-15T07:15:22.098Z","version":2},
      |        {"platform":"ubuntu","message":"Take a breath, let it go, walk away","msgid":"ID643","pid":7432,"level":"notice","dt":"2023-02-15T07:15:23.098Z","version":1},
      |        {"platform":"nginx","message":"There's a breach in the warp core, captain","msgid":"ID57","pid":7010,"level":"crit","dt":"2023-02-15T07:15:24.099Z","version":2},
      |        {"platform":"ubuntu","message":"Great Scott! We're never gonna reach 88 mph with the flux capacitor in its current state!","msgid":"ID89","pid":1572,"level":"emerg","dt":"2023-02-15T07:15:25.098Z","version":1},
      |        {"platform":"kubernetes","message":"A bug was encountered but not in Vector, which doesn't have bugs","msgid":"ID246","pid":9083,"level":"err","dt":"2023-02-15T07:15:26.098Z","version":1},
      |        {"platform":"aws","message":"Great Scott! We're never gonna reach 88 mph with the flux capacitor in its current state!","msgid":"ID199","pid":1638,"level":"err","dt":"2023-02-15T07:15:27.098Z","version":2},
      |        {"platform":"kubernetes","message":"Pretty pretty pretty good","msgid":"ID121","pid":9957,"level":"emerg","dt":"2023-02-15T07:15:28.098Z","version":1},
      |        {"platform":"aws","message":"Maybe we just shouldn't use computers","msgid":"ID320","pid":3790,"level":"warning","dt":"2023-02-15T07:15:29.098Z","version":1},
      |        {"platform":"ubuntu","message":"A bug was encountered but not in Vector, which doesn't have bugs","msgid":"ID210","pid":2863,"level":"debug","dt":"2023-02-15T07:15:30.098Z","version":1},
      |        {"platform":"kubernetes","message":"Take a breath, let it go, walk away","msgid":"ID605","pid":2503,"level":"warning","dt":"2023-02-15T07:15:31.098Z","version":1},
      |        {"platform":"kubernetes","message":"Take a breath, let it go, walk away","msgid":"ID825","pid":9861,"level":"info","dt":"2023-02-15T07:15:32.098Z","version":1},
      |        {"platform":"docker","message":"You're not gonna believe what just happened","msgid":"ID467","pid":7064,"level":"debug","dt":"2023-02-15T07:15:33.098Z","version":1},
      |        {"platform":"aws","message":"Great Scott! We're never gonna reach 88 mph with the flux capacitor in its current state!","msgid":"ID815","pid":6729,"level":"crit","dt":"2023-02-15T07:15:34.099Z","version":2},
      |        {"platform":"ubuntu","message":"We're gonna need a bigger boat","msgid":"ID573","pid":1070,"level":"debug","dt":"2023-02-15T07:15:35.098Z","version":2},
      |        {"platform":"nginx","message":"Pretty pretty pretty good","msgid":"ID676","pid":1143,"level":"debug","dt":"2023-02-15T07:15:36.098Z","version":2},
      |        {"platform":"gcp","message":"#hugops to everyone who has to deal with this","msgid":"ID660","pid":7827,"level":"debug","dt":"2023-02-15T07:15:37.098Z","version":2},
      |        {"platform":"aws","message":"Maybe we just shouldn't use computers","msgid":"ID749","pid":3096,"level":"info","dt":"2023-02-15T07:15:38.099Z","version":2},
      |        {"platform":"aws","message":"#hugops to everyone who has to deal with this","msgid":"ID490","pid":3559,"level":"err","dt":"2023-02-15T07:15:39.099Z","version":1},
      |        {"platform":"kubernetes","message":"Great Scott! We're never gonna reach 88 mph with the flux capacitor in its current state!","msgid":"ID76","pid":597,"level":"warning","dt":"2023-02-15T07:15:40.098Z","version":1},
      |        {"platform":"aws","message":"A bug was encountered but not in Vector, which doesn't have bugs","msgid":"ID265","pid":2323,"level":"crit","dt":"2023-02-15T07:15:41.098Z","version":2},
      |        {"platform":"kubernetes","message":"#hugops to everyone who has to deal with this","msgid":"ID449","pid":7798,"level":"warning","dt":"2023-02-15T07:15:42.098Z","version":1},
      |        {"platform":"aws","message":"Maybe we just shouldn't use computers","msgid":"ID13","pid":1168,"level":"emerg","dt":"2023-02-15T07:15:43.097Z","version":1},
      |        {"platform":"nginx","message":"Take a breath, let it go, walk away","msgid":"ID234","pid":3418,"level":"alert","dt":"2023-02-15T07:15:44.098Z","version":1},
      |        {"platform":"ubuntu","message":"There's a breach in the warp core, captain","msgid":"ID280","pid":8084,"level":"err","dt":"2023-02-15T07:15:45.097Z","version":1},
      |        {"platform":"aws","message":"Pretty pretty pretty good","msgid":"ID883","pid":8071,"level":"crit","dt":"2023-02-15T07:15:46.098Z","version":2},
      |        {"platform":"ubuntu","message":"We're gonna need a bigger boat","msgid":"ID489","pid":302,"level":"crit","dt":"2023-02-15T07:15:47.097Z","version":2},
      |        {"platform":"kubernetes","message":"We're gonna need a bigger boat","msgid":"ID264","pid":8421,"level":"crit","dt":"2023-02-15T07:15:48.097Z","version":1},
      |        {"platform":"kubernetes","message":"A bug was encountered but not in Vector, which doesn't have bugs","msgid":"ID218","pid":139,"level":"err","dt":"2023-02-15T07:15:49.098Z","version":2},
      |        {"platform":"ubuntu","message":"Great Scott! We're never gonna reach 88 mph with the flux capacitor in its current state!","msgid":"ID400","pid":9760,"level":"notice","dt":"2023-02-15T07:15:50.098Z","version":1},
      |        {"platform":"heroku","message":"There's a breach in the warp core, captain","msgid":"ID16","pid":1954,"level":"warning","dt":"2023-02-15T07:15:51.098Z","version":2},
      |        {"platform":"aws","message":"Maybe we just shouldn't use computers","msgid":"ID382","pid":3511,"level":"emerg","dt":"2023-02-15T07:15:52.098Z","version":1},
      |        {"platform":"gcp","message":"Pretty pretty pretty good","msgid":"ID919","pid":2460,"level":"notice","dt":"2023-02-15T07:15:53.099Z","version":1},
      |        {"platform":"gcp","message":"We're gonna need a bigger boat","msgid":"ID593","pid":2894,"level":"crit","dt":"2023-02-15T07:15:54.098Z","version":1},
      |        {"platform":"aws","message":"There's a breach in the warp core, captain","msgid":"ID141","pid":683,"level":"notice","dt":"2023-02-15T07:15:55.099Z","version":1},
      |        {"platform":"kubernetes","message":"Take a breath, let it go, walk away","msgid":"ID881","pid":8054,"level":"err","dt":"2023-02-15T07:15:56.100Z","version":2},
      |        {"platform":"nginx","message":"We're gonna need a bigger boat","msgid":"ID101","pid":4686,"level":"alert","dt":"2023-02-15T07:15:57.098Z","version":2},
      |        {"platform":"kubernetes","message":"Maybe we just shouldn't use computers","msgid":"ID771","pid":4775,"level":"info","dt":"2023-02-15T07:15:58.098Z","version":2},
      |        {"platform":"nginx","message":"A bug was encountered but not in Vector, which doesn't have bugs","msgid":"ID509","pid":5825,"level":"notice","dt":"2023-02-15T07:15:59.098Z","version":1},
      |        {"platform":"heroku","message":"There's a breach in the warp core, captain","msgid":"ID318","pid":3382,"level":"alert","dt":"2023-02-15T07:16:00.097Z","version":2},
      |        {"platform":"aws","message":"Take a breath, let it go, walk away","msgid":"ID96","pid":6858,"level":"notice","dt":"2023-02-15T07:16:01.097Z","version":2},
      |        {"platform":"kubernetes","message":"You're not gonna believe what just happened","msgid":"ID1","pid":4149,"level":"crit","dt":"2023-02-15T07:16:02.098Z","version":2},
      |        {"platform":"aws","message":"Pretty pretty pretty good","msgid":"ID86","pid":8238,"level":"notice","dt":"2023-02-15T07:16:03.098Z","version":2},
      |        {"platform":"gcp","message":"Maybe we just shouldn't use computers","msgid":"ID609","pid":948,"level":"warning","dt":"2023-02-15T07:16:04.098Z","version":2},
      |        {"platform":"docker","message":"A bug was encountered but not in Vector, which doesn't have bugs","msgid":"ID89","pid":7821,"level":"crit","dt":"2023-02-15T07:16:05.097Z","version":1},
      |        {"platform":"kubernetes","message":"Great Scott! We're never gonna reach 88 mph with the flux capacitor in its current state!","msgid":"ID600","pid":7243,"level":"notice","dt":"2023-02-15T07:16:06.098Z","version":1},
      |        {"platform":"heroku","message":"You're not gonna believe what just happened","msgid":"ID154","pid":2668,"level":"info","dt":"2023-02-15T07:16:07.097Z","version":1},
      |        {"platform":"kubernetes","message":"A bug was encountered but not in Vector, which doesn't have bugs","msgid":"ID22","pid":6477,"level":"err","dt":"2023-02-15T07:16:08.098Z","version":1},
      |        {"platform":"docker","message":"Great Scott! We're never gonna reach 88 mph with the flux capacitor in its current state!","msgid":"ID265","pid":8778,"level":"debug","dt":"2023-02-15T07:16:09.097Z","version":1},
      |        {"platform":"heroku","message":"Pretty pretty pretty good","msgid":"ID129","pid":6173,"level":"emerg","dt":"2023-02-15T07:16:10.098Z","version":1},
      |        {"platform":"kubernetes","message":"Maybe we just shouldn't use computers","msgid":"ID820","pid":2166,"level":"emerg","dt":"2023-02-15T07:16:11.098Z","version":2},
      |        {"platform":"aws","message":"A bug was encountered but not in Vector, which doesn't have bugs","msgid":"ID839","pid":319,"level":"crit","dt":"2023-02-15T07:16:12.098Z","version":2},
      |        {"platform":"nginx","message":"Maybe we just shouldn't use computers","msgid":"ID26","pid":4514,"level":"debug","dt":"2023-02-15T07:16:13.098Z","version":1},
      |        {"platform":"heroku","message":"Pretty pretty pretty good","msgid":"ID451","pid":7799,"level":"info","dt":"2023-02-15T07:16:14.097Z","version":1},
      |        {"platform":"aws","message":"A bug was encountered but not in Vector, which doesn't have bugs","msgid":"ID954","pid":9481,"level":"err","dt":"2023-02-15T07:16:15.098Z","version":1},
      |        {"platform":"heroku","message":"You're not gonna believe what just happened","msgid":"ID946","pid":2703,"level":"crit","dt":"2023-02-15T07:16:16.097Z","version":1},
      |        {"platform":"docker","message":"Take a breath, let it go, walk away","msgid":"ID419","pid":6899,"level":"info","dt":"2023-02-15T07:16:17.097Z","version":1},
      |        {"platform":"gcp","message":"There's a breach in the warp core, captain","msgid":"ID893","pid":9026,"level":"info","dt":"2023-02-15T07:16:18.098Z","version":1},
      |        {"platform":"aws","message":"We're gonna need a bigger boat","msgid":"ID697","pid":7542,"level":"warning","dt":"2023-02-15T07:16:19.098Z","version":2},
      |        {"platform":"gcp","message":"We're gonna need a bigger boat","msgid":"ID747","pid":983,"level":"crit","dt":"2023-02-15T07:16:20.098Z","version":2},
      |        {"platform":"docker","message":"We're gonna need a bigger boat","msgid":"ID587","pid":2321,"level":"crit","dt":"2023-02-15T07:16:21.098Z","version":2},
      |        {"platform":"ubuntu","message":"Pretty pretty pretty good","msgid":"ID447","pid":5505,"level":"crit","dt":"2023-02-15T07:16:22.097Z","version":2},
      |        {"platform":"ubuntu","message":"Maybe we just shouldn't use computers","msgid":"ID109","pid":7678,"level":"notice","dt":"2023-02-15T07:16:23.098Z","version":1},
      |        {"platform":"gcp","message":"There's a breach in the warp core, captain","msgid":"ID74","pid":5729,"level":"warning","dt":"2023-02-15T07:16:24.098Z","version":2},
      |        {"platform":"ubuntu","message":"You're not gonna believe what just happened","msgid":"ID959","pid":9238,"level":"emerg","dt":"2023-02-15T07:16:25.097Z","version":1},
      |        {"platform":"aws","message":"#hugops to everyone who has to deal with this","msgid":"ID423","pid":3782,"level":"warning","dt":"2023-02-15T07:16:26.098Z","version":1},
      |        {"platform":"heroku","message":"Maybe we just shouldn't use computers","msgid":"ID570","pid":5196,"level":"notice","dt":"2023-02-15T07:16:27.098Z","version":2},
      |        {"platform":"heroku","message":"There's a breach in the warp core, captain","msgid":"ID747","pid":2391,"level":"warning","dt":"2023-02-15T07:16:28.097Z","version":2},
      |        {"platform":"kubernetes","message":"We're gonna need a bigger boat","msgid":"ID283","pid":7671,"level":"debug","dt":"2023-02-15T07:16:29.097Z","version":2},
      |        {"platform":"ubuntu","message":"Great Scott! We're never gonna reach 88 mph with the flux capacitor in its current state!","msgid":"ID77","pid":1755,"level":"warning","dt":"2023-02-15T07:16:30.098Z","version":1},
      |        {"platform":"gcp","message":"There's a breach in the warp core, captain","msgid":"ID62","pid":9763,"level":"notice","dt":"2023-02-15T07:16:31.097Z","version":1},
      |        {"platform":"gcp","message":"Take a breath, let it go, walk away","msgid":"ID87","pid":568,"level":"debug","dt":"2023-02-15T07:16:32.098Z","version":2},
      |        {"platform":"heroku","message":"A bug was encountered but not in Vector, which doesn't have bugs","msgid":"ID191","pid":7848,"level":"warning","dt":"2023-02-15T07:16:33.098Z","version":2},
      |        {"platform":"gcp","message":"Great Scott! We're never gonna reach 88 mph with the flux capacitor in its current state!","msgid":"ID178","pid":5907,"level":"crit","dt":"2023-02-15T07:16:34.097Z","version":1},
      |        {"platform":"nginx","message":"A bug was encountered but not in Vector, which doesn't have bugs","msgid":"ID715","pid":2429,"level":"alert","dt":"2023-02-15T07:16:35.098Z","version":1}
      |    ]
      |}
      |""".stripMargin


    val inputSchemaStr =
      """
        |{
        |  "type": "record",
        |  "name": "record",
        |  "fields": [
        |    {
        |      "name": "platform",
        |      "type": "string"
        |    },
        |    {
        |      "name": "message",
        |      "type": "string"
        |    },
        |    {
        |      "name": "msgid",
        |      "type": "string"
        |    },
        |    {
        |      "name": "pid",
        |      "type": "int"
        |    },
        |    {
        |      "name": "level",
        |      "type": "string"
        |    },
        |    {
        |      "name": "dt",
        |      "type": "string"
        |    },
        |    {
        |      "name": "version",
        |      "type": "int"
        |    }
        |  ]
        |}""".stripMargin

    val inputSchema = new Schema.Parser().parse(inputSchemaStr)


    val decoded: Either[Error, Data] = decode[Data](jsonRecordsString)

    decoded.right.get.data.map(r => {
      val inputRecord = new GenericData.Record(inputSchema)
      inputRecord.put("platform", r.platform)
      inputRecord.put("message", r.message)
      inputRecord.put("msgid", r.msgid)
      inputRecord.put("pid", r.pid)
      inputRecord.put("level", r.level)
      inputRecord.put("dt", r.dt)
      inputRecord.put("version", r.version)
      inputRecord
    })
  }

}
