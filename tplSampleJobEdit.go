package main

import (
	"fmt"
	"html"

	"github.com/kr/beanstalk"
)

// tplSampleJobEdit render a sample job edit form.
func tplSampleJobEdit(key string, alert string) string {
	var err error
	var action, title, name, savedTo, saveTo, data, ST string
	if key == "" {
		action = `?action=actionNewSample`
		title = `<h4 class="text-info">New sample job</h4>`
	} else {
		action = `?action=actionEditSample&key=` + key
		name = html.EscapeString(getSampleJobNameByKey(key))
		data = html.EscapeString(getSampleJobDataByKey(key))
		title = fmt.Sprintf(`<h4 class="text-info">Edit: %s</h4>`, name)
		for _, j := range sampleJobs.Jobs {
			if key == j.Key {
				for _, t := range j.Tubes {
					saveTo += fmt.Sprintf(`<div class="control-group">
                                <div class="controls">
                                    <label class="checkbox-inline">
                                        <input type="checkbox" name="tubes[%s]" value="1" checked="checked">
                                        %s
                                    </label>
                                </div>
                            </div>`, t, t)
				}
			}
		}
	}

	for _, server := range selfConf.Servers {
		var bstkConn *beanstalk.Conn
		var tubeList string
		if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
			continue
		}
		tubes, _ := bstkConn.ListTubes()
		bstkConn.Close()
		for _, v := range tubes {
			var checked string
			for _, j := range sampleJobs.Jobs {
				if j.Key == key {
					for _, t := range j.Tubes {
						if t == v {
							checked = `checked="checked"`
						}
					}
				}
			}
			tubeList += fmt.Sprintf(`<div class="control-group">
                                        <div class="controls">
                                            <label class="checkbox-inline">
                                                <input type="checkbox" name="tubes[%s]" value="1" %s>
                                                %s
                                            </label>
                                        </div>
                                    </div>`, v, checked, v)
		}
		ST += fmt.Sprintf(`<div class="pull-left" style="padding-right: 35px;">
                            %s
                            <blockquote>
                                %s
                            </blockquote>
                        </div>`, server, tubeList)
	}
	if name != "" {
		savedTo = fmt.Sprintf(`<div class="pull-left" style="padding-right: 35px;">
                Saved to:
                <blockquote>
                    %s
                </blockquote>
            </div>`, saveTo)
	}

	return fmt.Sprintf(`<form name="sampleJobsEdit" action="%s" method="POST">
    <div class="clearfix form-group">
        <div class="pull-left">
            %s
        </div>
        <div class="pull-right">
            <a href="?action=manageSamples" class="btn btn-default btn-small"><i class="glyphicon glyphicon-list"></i> Manage samples</a>
        </div>
    </div>
    <div class=" form-group">
        <fieldset>
            %s
            <div class="control-group">
                <label class="control-label" for="addsamplename"><b>Name *</b></label>

                <div class="controls form-group">
                    <input class="input-xlarge focused" id="addsamplename" name="name" type="text" value="%s"
                           autocomplete="off">
                </div>
            </div>
        </fieldset>
        <div class="clearfix">
            <label class="control-label"><b>Available on tubes *</b></label>
            <br/>
            %s
            %s
        </div>
        <div>
            <label class="control-label" for="jobdata"><b>Job data *</b></label>
            <textarea name="jobdata" id="jobdata" %s</textarea>
        </div>
    </div>
    <div>
        <input type="submit" class="btn btn-success" value="Save"/>
    </div>
</form>`, action, title, alert, name, savedTo, ST, `style="width:100%" rows="3">`+data)
}
