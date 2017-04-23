from lsbucketbeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Lsbucketbeat normally
        """
        self.render_config_template(
            path=os.path.abspath(self.working_dir) + "/log/*"
        )

        lsbucketbeat_proc = self.start_beat()
        self.wait_until(lambda: self.log_contains("lsbucketbeat is running"))
        exit_code = lsbucketbeat_proc.kill_and_wait()
        assert exit_code == 0
